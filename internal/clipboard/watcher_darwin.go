// internal/clipboard/watcher_darwin.go
package clipboard

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit
#import <AppKit/AppKit.h>
#import <stdlib.h>

static int pasteboard_change_count(void) {
    return (int)[[NSPasteboard generalPasteboard] changeCount];
}

typedef struct {
    char*          uti;
    char*          text;
    unsigned char* data;
    int            data_len;
} PasteboardData;

static PasteboardData read_pasteboard(void) {
    PasteboardData pd = {NULL, NULL, NULL, 0};
    NSPasteboard *pb  = [NSPasteboard generalPasteboard];

    // Binary image types (priority order)
    for (NSString *t in @[NSPasteboardTypePNG, NSPasteboardTypeTIFF,
                          @"com.compuserve.gif", @"public.heic"]) {
        NSData *d = [pb dataForType:t];
        if (d) {
            pd.uti  = strdup(t.UTF8String);
            pd.data = (unsigned char*)malloc(d.length);
            memcpy(pd.data, d.bytes, d.length);
            pd.data_len = (int)d.length;
            return pd;
        }
    }

    // NSColor (from color pickers)
    if ([[pb types] containsObject:NSPasteboardTypeColor]) {
        NSColor *c = [[NSColor colorFromPasteboard:pb]
                       colorUsingColorSpace:[NSColorSpace sRGBColorSpace]];
        if (c) {
            pd.uti = strdup("com.apple.cocoa.pasteboard.color");
            char *hex = (char*)malloc(8);
            snprintf(hex, 8, "#%02x%02x%02x",
                (int)(c.redComponent   * 255),
                (int)(c.greenComponent * 255),
                (int)(c.blueComponent  * 255));
            pd.text = hex;
            return pd;
        }
    }

    // File URL — skip directories
    if ([[pb types] containsObject:@"public.file-url"]) {
        NSString *urlStr = [pb stringForType:@"public.file-url"];
        NSURL    *url    = [NSURL URLWithString:urlStr];
        NSNumber *isDir  = nil;
        [url getResourceValue:&isDir forKey:NSURLIsDirectoryKey error:nil];
        if (![isDir boolValue]) {
            pd.uti  = strdup("public.file-url");
            pd.text = strdup(url.path.UTF8String);
            return pd;
        }
    }

    // RTF → extract plain text
    NSData *rtf = [pb dataForType:NSPasteboardTypeRTF];
    if (rtf) {
        NSAttributedString *as = [[NSAttributedString alloc]
            initWithRTF:rtf documentAttributes:nil];
        pd.uti  = strdup("public.rtf");
        pd.text = strdup(as.string.UTF8String ?: "");
        return pd;
    }

    // Plain text (checked before HTML — most apps write both, prefer plain)
    NSString *str = [pb stringForType:NSPasteboardTypeString];
    if (str) {
        pd.uti  = strdup("public.utf8-plain-text");
        pd.text = strdup(str.UTF8String);
        return pd;
    }

    // HTML (only when no plain text is available)
    NSString *html = [pb stringForType:NSPasteboardTypeHTML];
    if (html) {
        pd.uti  = strdup("public.html");
        pd.text = strdup(html.UTF8String);
        return pd;
    }

    return pd;
}

static void free_pasteboard_data(PasteboardData pd) {
    if (pd.uti)  free(pd.uti);
    if (pd.text) free(pd.text);
    if (pd.data) free(pd.data);
}

static void write_text_to_pasteboard(const char* text) {
    NSString *s = [NSString stringWithUTF8String:text];
    [[NSPasteboard generalPasteboard] clearContents];
    [[NSPasteboard generalPasteboard] setString:s forType:NSPasteboardTypeString];
}

static void write_image_to_pasteboard(const char* path) {
    NSImage *img = [[NSImage alloc] initWithContentsOfFile:
                    [NSString stringWithUTF8String:path]];
    if (!img) return;
    [[NSPasteboard generalPasteboard] clearContents];
    [[NSPasteboard generalPasteboard] writeObjects:@[img]];
}
*/
import "C"
import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
	"unsafe"
)

// RawItem is a clipboard capture before classification or DB write.
type RawItem struct {
	UTI  string // NSPasteboard type identifier
	Text string // populated for text-based types
	Data []byte // populated for binary types (images)
}

// TextHash returns a SHA256 hex string for text content.
func TextHash(text string) string {
	h := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", h)
}

// Watcher polls NSPasteboard.changeCount on an adaptive interval.
type Watcher struct {
	pause  chan struct{}
	resume chan struct{}
}

func NewWatcher() *Watcher {
	return &Watcher{
		pause:  make(chan struct{}, 1),
		resume: make(chan struct{}, 1),
	}
}

// Pause stops polling (call on screen lock).
func (w *Watcher) Pause() { w.pause <- struct{}{} }

// Resume restarts polling (call on screen unlock).
func (w *Watcher) Resume() { w.resume <- struct{}{} }

// Start blocks until ctx is cancelled. Each new clipboard change sends a RawItem to out.
func (w *Watcher) Start(ctx context.Context, out chan<- RawItem) {
	lastCount    := int(C.pasteboard_change_count())
	lastActivity := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.pause:
			<-w.resume // block until resumed
			continue
		default:
		}

		interval := 250 * time.Millisecond
		if time.Since(lastActivity) > 5*time.Minute {
			interval = 2 * time.Second
		}
		time.Sleep(interval)

		current := int(C.pasteboard_change_count())
		if current == lastCount {
			continue
		}
		lastCount = current
		lastActivity = time.Now()

		pd := C.read_pasteboard()
		if pd.uti == nil {
			continue
		}

		raw := RawItem{UTI: C.GoString(pd.uti)}
		if pd.text != nil {
			raw.Text = C.GoString(pd.text)
		}
		if pd.data != nil && pd.data_len > 0 {
			raw.Data = C.GoBytes(unsafe.Pointer(pd.data), pd.data_len)
		}
		C.free_pasteboard_data(pd)

		select {
		case out <- raw:
		case <-ctx.Done():
			return
		}
	}
}

// WriteTextToClipboard puts text back onto NSPasteboard.
func WriteTextToClipboard(text string) {
	cs := C.CString(text)
	defer C.free(unsafe.Pointer(cs))
	C.write_text_to_pasteboard(cs)
}

// WriteImageToClipboard puts an image file back onto NSPasteboard.
func WriteImageToClipboard(path string) {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))
	C.write_image_to_pasteboard(cs)
}

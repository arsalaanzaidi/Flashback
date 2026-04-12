export namespace store {
	
	export class Item {
	    id: string;
	    content: string;
	    contentHash: string;
	    type: string;
	    subtype: string;
	    pinned: boolean;
	    copiedAt: number;
	    createdAt: number;
	    charCount: number;
	    imagePath: string;
	    thumbBase64: string;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.content = source["content"];
	        this.contentHash = source["contentHash"];
	        this.type = source["type"];
	        this.subtype = source["subtype"];
	        this.pinned = source["pinned"];
	        this.copiedAt = source["copiedAt"];
	        this.createdAt = source["createdAt"];
	        this.charCount = source["charCount"];
	        this.imagePath = source["imagePath"];
	        this.thumbBase64 = source["thumbBase64"];
	    }
	}
	export class Settings {
	    retentionMode: string;
	    retentionValue: number;
	    globalShortcut: string;
	    launchAtLogin: boolean;
	    firstRunComplete: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.retentionMode = source["retentionMode"];
	        this.retentionValue = source["retentionValue"];
	        this.globalShortcut = source["globalShortcut"];
	        this.launchAtLogin = source["launchAtLogin"];
	        this.firstRunComplete = source["firstRunComplete"];
	    }
	}

}


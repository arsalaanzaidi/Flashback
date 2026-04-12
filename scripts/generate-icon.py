#!/usr/bin/env python3
"""Generate build/appicon.png — B5 style (1024×1024).

Design: dark navy bg, blue clipboard with gradient, 3 left speed lines,
gold bolt badge top-right.
"""
import os

from PIL import Image, ImageDraw

S = 1024
OUT = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'build', 'appicon.png')


def rgb(h):
    h = h.lstrip('#')
    return tuple(int(h[i:i+2], 16) for i in (0, 2, 4))


def rgba(h, a=255):
    return (*rgb(h), a)


img = Image.new('RGBA', (S, S), (0, 0, 0, 0))
d = ImageDraw.Draw(img)

# ── 1. Background (dark navy, rounded square) ──────────────────────────────────
d.rounded_rectangle([0, 0, S - 1, S - 1], radius=int(S * 0.22), fill=rgba('0f172a'))

# ── 2. Clipboard body (blue gradient, top=light, bottom=dark) ─────────────────
cbx1, cby1 = int(S * 0.27), int(S * 0.22)
cbx2, cby2 = int(S * 0.77), int(S * 0.87)
cb_r = int(S * 0.055)

grad = Image.new('RGBA', (S, S), (0, 0, 0, 0))
gd = ImageDraw.Draw(grad)
top_c, bot_c = rgb('60a5fa'), rgb('2563eb')
span = cby2 - cby1
for y in range(cby1, cby2 + 1):
    t = (y - cby1) / span
    col = tuple(int(top_c[i] + t * (bot_c[i] - top_c[i])) for i in range(3)) + (255,)
    gd.line([(cbx1, y), (cbx2, y)], fill=col)

mask = Image.new('L', (S, S), 0)
ImageDraw.Draw(mask).rounded_rectangle([cbx1, cby1, cbx2, cby2], radius=cb_r, fill=255)
img.paste(grad, mask=mask)
d = ImageDraw.Draw(img)  # refresh after paste

# ── 3. Clipboard clip (top centre piece) ──────────────────────────────────────
cl_w = int(S * 0.20)
cl_cx = (cbx1 + cbx2) // 2
cl_x1 = cl_cx - cl_w // 2
cl_y1 = cby1 - int(S * 0.04)
cl_y2 = cl_y1 + int(S * 0.09)
d.rounded_rectangle([cl_x1, cl_y1, cl_x1 + cl_w, cl_y2],
                    radius=int(S * 0.035), fill=rgba('93c5fd'))

inner_w = int(cl_w * 0.55)
ix1 = cl_cx - inner_w // 2
d.rounded_rectangle([ix1, cl_y1 - int(S * 0.01), ix1 + inner_w, cl_y1 + int(S * 0.045)],
                    radius=int(S * 0.022), fill=rgba('bfdbfe'))

# ── 4. Speed lines (left side, 3 lines, tapering width + opacity) ─────────────
lx0 = int(S * 0.04)
lx1 = int(S * 0.24)
base_w = max(1, int(S / 120 * 2.5))
speed_lines = [
    (int(S * 0.41), base_w,          int(255 * 0.50)),
    (int(S * 0.51), max(1, base_w-1), int(255 * 0.35)),
    (int(S * 0.61), max(1, base_w-2), int(255 * 0.22)),
]
for ly, lw, la in speed_lines:
    d.line([(lx0, ly), (lx1, ly)], fill=(*rgb('3b82f6'), la), width=lw)

# ── 5. Content lines on clipboard face ────────────────────────────────────────
lh = int(S * 0.015)
lr = int(S * 0.012)
lx_inner = int(S * 0.335)
content_lines = [
    (int(S * 0.43), int(S * 0.265), int(255 * 0.50)),
    (int(S * 0.52), int(S * 0.195), int(255 * 0.32)),
    (int(S * 0.61), int(S * 0.230), int(255 * 0.32)),
    (int(S * 0.70), int(S * 0.145), int(255 * 0.20)),
]
for ly, lw, la in content_lines:
    d.rounded_rectangle([lx_inner, ly - lh, lx_inner + lw, ly + lh],
                        radius=lr, fill=(255, 255, 255, la))

# ── 6. Badge: gold circle, top-right ──────────────────────────────────────────
bcx, bcy = int(S * 0.775), int(S * 0.27)
br = int(S * 0.115)
d.ellipse([bcx - br, bcy - br, bcx + br, bcy + br], fill=rgba('fbbf24'))

# ── 7. Lightning bolt inside badge (7-point white polygon) ────────────────────
# Proportions derived from reference SVG (bolt in circle of radius 14):
#   points="92,25 86,38 91,38 84,48 98,34 92,34 96,25" center=(92,34)
# Scaled to unit radius then multiplied by br.
bolt_rel = [
    ( 0.00, -0.64),  # top
    (-0.43,  0.29),  # mid-left
    (-0.07,  0.29),  # mid inner
    (-0.57,  1.00),  # bottom tip
    ( 0.43,  0.00),  # right outer
    ( 0.00,  0.00),  # centre
    ( 0.29, -0.64),  # upper-right
]
bolt_pts = [(bcx + int(bx * br), bcy + int(by * br)) for bx, by in bolt_rel]
d.polygon(bolt_pts, fill=(255, 255, 255, 255))

# ── Save ───────────────────────────────────────────────────────────────────────
img.save(OUT, 'PNG')
print(f"Saved: {OUT}  ({S}x{S})")

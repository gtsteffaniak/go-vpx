# Rejected subpel MV-rate floor probe

Precompute a true lower bound while building the existing component-rate LUT. Stop subpixel search at a center whose variance-plus-rate reaches that floor. Ordinary stream and frozen stats-walk results match.

All-search probe: fixed codec paired 0.95156/0.95344x; adaptive 1.07804/1.02016x. Narrow iterative-only follow-up: fixed 0.97466/0.98791x, adaptive 0.97566/1.00593x. First follow-up pair had competing compilation. No robust whole-codec gain; keep out of production. Source/field layout and additional setup may matter; not established by this screen.

Raw codec samples and unsuccessful variants are retained.

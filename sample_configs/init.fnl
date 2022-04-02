;; Default fnl Hilbish config
(local colors (require :lunacolors))
(local b (require :bait))
(local ansikit (require :ansikit))

(fn doPrompt [fail]
    (hilbish.prompt (colors.format
                     (.. "{blue}%u {cyan}%d " (if (= fail 0) "{green}" "{red}") "∆ "))))
(doPrompt 0)
(print (colors.format hilbish.greeting))

(b.catch "command.exit" doPrompt)

(b.catch "hilbish.vimMode" (lambda [mode]
                             (if (= mode "insert")
                                 (ansikit.cursorStyle ansikit.lineCursor)
                                 (ansikit.cursorStyle ansikit.blockCursor))
                             (values)))

# Linux Man Pages in Dash
Oct 21st, 2013

Dash works great as a man page browser, but I sometimes get requests to make extra docsets containing the man pages of various flavours of Linux.

I’ve decided not to pursue these requests, because:

Updates for these docsets would be a nightmare, as man pages change a lot, individually.
I’d have to choose which man pages to include and which not to. I’d never be able to guess which obscure man page a user might want.
The current Man Pages docset solves these issues by indexing the man pages that are actually on your Mac.

The workaround
You can copy the man pages from any Linux box to your Mac and Dash will index them as part of the regular Man Pages docset.

Step by step instructions:

Log into your Linux box
Run man -w to list the folders that contain man pages
Copy these folders to your Mac
Optional, but highly recommended: use a batch renamer to rename all of the man page files to have a common prefix. This will help you differentiate between the default macOS man pages and the Linux ones
Move the man pages anywhere on your MANPATH, or any folder from man -w
That’s it!
# iTerm2 下使用 rz/sz 命令

> 参考：[ZModem integration for iTerm 2](https://github.com/mmastrac/iterm2-zmodem)


步骤：

- Install `lrzsz` on OSX: `brew install lrzsz`
- Save the `iterm2-send-zmodem.sh` and `iterm2-recv-zmodem.sh` scripts in `/usr/local/bin/`
- Set up Triggers (Perferences -> Profiles -> your Profile -> Advanced -> Triggers -> Edit) in iTerm 2 like so:
```
    Regular expression: rz waiting to receive.\*\*B0100
    Action: Run Silent Coprocess
    Parameters: /usr/local/bin/iterm2-send-zmodem.sh
    Instant: checked

    Regular expression: \*\*B00000000000000
    Action: Run Silent Coprocess
    Parameters: /usr/local/bin/iterm2-recv-zmodem.sh
    Instant: checked
```
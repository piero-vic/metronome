# metronome

A small metronome for the terminal made with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and
[Beep](https://github.com/gopxl/beep).

## Installation

Clone this repository and run:

```bash
go build
```

or

```bash
go install
```

## Usage

### Flags

```
Usage: metronome [options]

Options:
  -t, --tempo=[int]   set the tempo in BPM (default 60)
  -b, --beats=[int]   set the number of beats per measure (default 4)
  -p, --play          start the metronome automatically
  -h, --help          print help
```

### Keys

| Key                                                          | Description                          |
| ------------------------------------------------------------ | ------------------------------------ |
| <kbd>Up</kbd> <kbd>Down</kbd> or <kbd>k</kbd> <kbd>j</kbd>   | Increase or decrease tempo           |
| <kbd>Right</kbd><kbd>Left</kbd> or <kbd>l</kbd> <kbd>h</kbd> | Increase or decrease number of beats |
| <kbd>Space</kbd> or <kbd>p</kbd>                             | Play / Pause                         |
| <kbd>Ctrl + c</kbd> or <kbd>q</kbd>                          | Quit                                 |

## License

This project is licensed under the [MIT License](LICENSE).

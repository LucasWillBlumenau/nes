# NES

This repository contains the source code and the binaries for a Nintendo Entertainment System emulator written in Go. 

To use the app, you can run the following command on Linux:

```bash
./nes <path-to-your-rom-file>
```

and this one on Windows:


```powershell
.\nes.exe <path-to-your-rom-file>
```

## Playing The Games

In order to play the games, you can use a controller or the keyboard. In case, you use the keyboards this is the relationship between the keys and the buttons of the NES controller:

- W -> Button Up
- S -> Button Down
- A -> Button Left
- D -> Button Rigth
- Left Shift -> Button A
- Space -> Button B
- Backspace -> Select
- Enter -> Start


## Notes

The project is still in development, and a lot of games shouldn't be running yet, and some features might
be missing.

## Useful Resources

Here are some of the resources that I used to build the emulator, which were really useful:

- [Nesdev Wiki](https://www.nesdev.org/wiki/Nesdev_Wiki)
- [Pagetable 6502 Instruction Set](https://www.pagetable.com/c64ref/6502/?tab=2#)

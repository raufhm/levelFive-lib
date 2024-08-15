# Printer & Parser Lib

This Go package provides functionalities for:

- **Parsing Templates**: Format ticketing and POS print messages using templates.
- **Printing**: Print to USB printers or over the network.

## Features

- **Template Parsing**: Easily format and parse ticket and POS print messages.
- **Printer Detection**: Identify local USB printers and network printers on macOS.
- **Local and Remote Printing**: Print documents to USB printers or remote printers over SSH.

## Requirements

- **macOS**: Utilizes the `lpstat` command. Support for other operating systems may be added in the future.

## Installation

Clone the repository and import it into your Go project:

```sh
git clone https://github.com/yourusername/printer-utility.git

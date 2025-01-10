# go-tuner

## Overview
`go-tuner` is a Go-based audio tuner application that uses WASAPI for audio capture. It captures audio from a specified input device, processes the audio data using FFT (Fast Fourier Transform), and detects the dominant pitch frequency.

## Features
- Audio capture using WASAPI
- Real-time pitch detection
- Signal handling for graceful shutdown

## Requirements
- Go 1.23.1 or later
- A compatible audio input device

## Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/go-tuner.git
    cd go-tuner
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

## Usage
1. Run the application:
    ```sh
    go run main.go
    ```

2. The application will start capturing audio from the specified input device and print the detected pitch frequency to the console.

## Configuration
- The input device is currently hardcoded as `"Analogue 1 + 2 (Focusrite USB Audio)"`. You can change this in the `main.go` file.
- Sample Rate is also hardcoded to 48000 hz.

## License
This project is licensed under the MIT License.

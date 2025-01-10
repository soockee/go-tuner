package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"log/slog"
	"math/cmplx"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/soockee/go-record"
	"gonum.org/v1/gonum/dsp/fourier"
)

func main() {
	stream := record.NewAudioStream()
	ctx, _ := setupSignalHandling()

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		record.Capture(stream, "Analogue 1 + 2 (Focusrite USB Audio)", ctx)
	}()
	go func() {
		defer wg.Done()
		analyze(stream, ctx)
	}()
	wg.Wait()
}

func analyze(stream *record.AudioStream, ctx context.Context) {
	// format, err := record.GetFormat("Analogue 1 + 2 (Focusrite USB Audio)")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(15 * time.Millisecond)
			data := stream.Read()
			data = readLastBytes(data, 48000)
			buffer, err := bytesToFloat32(data)
			if err != nil {
				fmt.Println(err)
			}
			if err != nil {
				log.Printf("Error reading audio: %v", err)
				continue
			}

			pitch := processAudio(buffer) // Pass the buffer directly
			fmt.Printf("Detected pitch: %f Hz\n", pitch)
		}
	}
}

func readLastBytes(buffer []byte, n int) []byte {
	if len(buffer) < n {
		return buffer
	}
	return buffer[len(buffer)-n:]
}

func setupSignalHandling() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	// Catch OS signals like Ctrl+C
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		slog.Info("SIGINT received, shutting down...")
		cancel()
	}()

	return ctx, cancel
}

func processAudio(in []float32) float64 {
	// Convert float32 to float64 for FFT
	data := make([]float64, len(in))
	for i, v := range in {
		data[i] = float64(v)
	}
	if len(data) == 0 {
		return 0
	}

	// Create an FFT plan
	fft := fourier.NewFFT(len(data))
	// This performs the FFT and returns complex coefficients
	coeff := fft.Coefficients(nil, data)

	// Find dominant frequency
	return findDominantFrequency(coeff)
}

func findDominantFrequency(coeff []complex128) float64 {
	maxVal := 0.0
	var maxIdx int
	for i, v := range coeff {
		if abs := cmplx.Abs(v); abs > maxVal {
			maxVal = abs
			maxIdx = i
		}
	}
	sampleRate := 48000 // Define as per your setup
	// Calculate frequency
	return float64(maxIdx) * float64(sampleRate) / float64(len(coeff))
}

func bytesToFloat32(data []byte) ([]float32, error) {
	if len(data)%4 != 0 {
		return nil, fmt.Errorf("invalid byte slice length")
	}
	floatData := make([]float32, len(data)/4)
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &floatData)
	if err != nil {
		return nil, fmt.Errorf("failed to convert bytes to float32: %w", err)
	}
	return floatData, nil
}

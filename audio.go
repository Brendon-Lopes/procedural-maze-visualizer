package main

import (
	"encoding/binary"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

const carveWavPath = "carve.wav"

// 5 semitones down = 2^(5/12) ≈ 1.3348
const pitchShiftRatio = 1.3348399

// Subtle pitch variation per play: ±1% ≈ ±0.17 semitones
const pitchVariation = 0.01

// Fade in/out duration (in samples) to prevent clicks
const fadeSamples = sampleRate * 5 / 1000

var audioCtx *audio.Context

func init() {
	audioCtx = audio.NewContext(sampleRate)
}

type Sounds struct {
	carveBuf []byte
}

func NewSounds() *Sounds {
	carveData := loadWav(carveWavPath)

	if carveData == nil {
		return &Sounds{}
	}

	return &Sounds{carveBuf: carveData}
}

func (s *Sounds) PlayCarve() {
	if audioCtx == nil || s.carveBuf == nil {
		return
	}
	// Subtle pitch variation
	ratio := 1.0 + rand.Float64()*2*pitchVariation - pitchVariation
	buf := applyReverb(applyFade(pitchShift(s.carveBuf, ratio)))
	p := audioCtx.NewPlayerFromBytes(buf)
	p.Play()
}

func (s *Sounds) PlayBacktrack() {
	if audioCtx == nil || s.carveBuf == nil {
		return
	}
	// Subtle pitch variation around the base pitch shift
	ratio := pitchShiftRatio + rand.Float64()*2*pitchVariation - pitchVariation
	buf := applyReverb(applyFade(pitchShift(s.carveBuf, ratio)))
	p := audioCtx.NewPlayerFromBytes(buf)
	p.Play()
}

// pitchShift resamples 16-bit stereo PCM data.
// ratio < 1 → higher and shorter
// ratio > 1 → lower and longer
func pitchShift(pcm []byte, ratio float64) []byte {
	frameBytes := 4
	numFrames := len(pcm) / frameBytes
	newNumFrames := int(float64(numFrames) * ratio)
	result := make([]byte, newNumFrames*frameBytes)

	for i := 0; i < newNumFrames; i++ {
		srcPos := float64(i) / ratio
		srcFrame := int(srcPos)
		frac := srcPos - float64(srcFrame)

		srcOffset := srcFrame * frameBytes
		nextOffset := srcOffset + frameBytes

		if nextOffset+frameBytes-1 >= len(pcm) {
			remaining := len(pcm) - srcOffset
			if remaining > 0 && i*frameBytes+remaining <= len(result) {
				copy(result[i*frameBytes:], pcm[srcOffset:srcOffset+remaining])
			}
			break
		}

		// Interpolate each channel (16-bit LE)
		for ch := 0; ch < 2; ch++ {
			byteOffset := ch * 2

			s0 := int16(pcm[srcOffset+byteOffset]) | int16(pcm[srcOffset+byteOffset+1])<<8
			s1 := int16(pcm[nextOffset+byteOffset]) | int16(pcm[nextOffset+byteOffset+1])<<8

			interpolated := s0 + int16(float64(s1-s0)*frac)

			result[i*frameBytes+byteOffset] = byte(interpolated)
			result[i*frameBytes+byteOffset+1] = byte(interpolated >> 8)
		}
	}

	return result
}

// applyFade applies 5ms fade-in and fade-out to the PCM buffer
// to prevent clicks at the start and end of the sound.
func applyFade(pcm []byte) []byte {
	frameBytes := 4
	numFrames := len(pcm) / frameBytes

	if numFrames <= fadeSamples*2 {
		return pcm
	}

	result := make([]byte, len(pcm))
	copy(result, pcm)

	// Fade in
	for i := 0; i < fadeSamples && i < numFrames; i++ {
		gain := float32(i) / float32(fadeSamples)
		for ch := 0; ch < 2; ch++ {
			byteOffset := i*frameBytes + ch*2
			sample := int16(pcm[byteOffset]) | int16(pcm[byteOffset+1])<<8
			faded := int16(float32(sample) * gain)
			result[byteOffset] = byte(faded)
			result[byteOffset+1] = byte(faded >> 8)
		}
	}

	// Fade out
	for i := 0; i < fadeSamples && i < numFrames; i++ {
		frameIdx := numFrames - 1 - i
		gain := float32(i) / float32(fadeSamples)
		for ch := 0; ch < 2; ch++ {
			byteOffset := frameIdx*frameBytes + ch*2
			sample := int16(pcm[byteOffset]) | int16(pcm[byteOffset+1])<<8
			faded := int16(float32(sample) * gain)
			result[byteOffset] = byte(faded)
			result[byteOffset+1] = byte(faded >> 8)
		}
	}

	return result
}

// applyReverb adds a subtle reverb effect by mixing a delayed copy
// of the sound with reduced amplitude.
func applyReverb(pcm []byte) []byte {
	frameBytes := 4
	numFrames := len(pcm) / frameBytes

	// 35ms delay with 15% feedback
	delayFrames := sampleRate * 35 / 1000
	feedback := float32(0.15)

	result := make([]byte, len(pcm))
	copy(result, pcm)

	for i := delayFrames; i < numFrames; i++ {
		for ch := 0; ch < 2; ch++ {
			srcOffset := (i-delayFrames)*frameBytes + ch*2
			dstOffset := i*frameBytes + ch*2

			src := int16(pcm[srcOffset]) | int16(pcm[srcOffset+1])<<8
			dst := int16(result[dstOffset]) | int16(result[dstOffset+1])<<8

			mixed := dst + int16(float32(src)*feedback)
			// Clamp to prevent overflow
			if mixed > 32767 {
				mixed = 32767
			} else if mixed < -32768 {
				mixed = -32768
			}

			result[dstOffset] = byte(mixed)
			result[dstOffset+1] = byte(mixed >> 8)
		}
	}

	return result
}

// loadWav extracts PCM data from the "data" chunk of a WAV file.
// Expected format: PCM 16-bit, stereo, sample rate == audioCtx (44100 Hz).
func loadWav(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	pos := 12
	for pos+8 <= len(data) {
		chunkID := string(data[pos : pos+4])
		chunkSize := int(binary.LittleEndian.Uint32(data[pos+4 : pos+8]))
		if chunkID == "data" {
			start := pos + 8
			end := start + chunkSize
			if end > len(data) {
				end = len(data)
			}
			return data[start:end]
		}
		pos += 8 + chunkSize
		if chunkSize%2 == 1 {
			pos++
		}
	}
	return nil
}

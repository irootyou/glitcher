package main

import (
    "flag"
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/gif"
    "image/jpeg"
    "image/png"
    "image/color/palette"
    "log"
    "math/rand"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"

    _ "image/jpeg"
    _ "image/png"
)

// Define ANSI escape codes for colors
const (
    reset     = "\033[0m"
    pink      = "\033[38;5;206m"
    lightPink = "\033[38;5;218m"
    darkPink  = "\033[38;5;162m"
)

// Messages map for easy management of log and print messages
var messages = map[string]string{
    "startSpinner":     " 🌸 Let's get started!",
    "stopSpinner":      " 🌺 All done!",
    "fetchImageError":  " ❌ Oopsie! Couldn't fetch the image: %v",
    "saveImageError":   " ❌ Oh no! Couldn't save the image: %v",
    "unsupportedFormat":" ❌ Uh-oh! Unsupported image format: %v",
    "encodingError":    " ❌ Yikes! Error encoding the image: %v",
    "glitchComplete":   " ✅ Yay! Image glitching complete!",
    "blackImageError":  " ❌ Oh dear! The image is completely black. Aborting.",
    "inputRequired":    " ⚠️  Please provide an input image file path, cutie!",
    "usage": `Usage: glitcher [flags] input_image [output_image]
Flags:
  -glitch-intensity float  Intensity of the glitch effect (0.1-9.0) (default 5.0)
  -scan-lines              Apply scan lines glitch effect (default false)
  -pixel-sort              Apply pixel sort glitch effect (default false)
  -color-offset            Apply color offset glitch effect (default false)
  -seed int                Random seed for reproducibility (default random)
  -gif                     Create a GIF instead of a single image (default false)
  -frames int              Number of frames for the GIF (default 10)
  -delay int               Delay between frames in the GIF (default 10)
  -cycle int               Number of cycles of glitches to apply (default 1)
  -step int                Step size between frames (default 1)
`,
}

// ImageGlitcher holds the necessary image settings and data
type ImageGlitcher struct {
    ImgWidth, ImgHeight int
    Image               *image.RGBA
    Seed                int64
}

// FetchImage loads an image from a file path and initializes the ImageGlitcher fields
func (ig *ImageGlitcher) FetchImage(srcImg string) error {
    img, err := ig.OpenImage(srcImg)
    if err != nil {
        return err
    }
    ig.ImgWidth, ig.ImgHeight = img.Bounds().Dx(), img.Bounds().Dy()
    ig.Image = imageToRGBA(img)
    return nil
}

// OpenImage opens an image file and decodes it
func (ig *ImageGlitcher) OpenImage(imgPath string) (image.Image, error) {
    file, err := os.Open(imgPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    img, _, err := image.Decode(file)
    return img, err
}

// imageToRGBA converts an image to RGBA format
func imageToRGBA(img image.Image) *image.RGBA {
    rgba := image.NewRGBA(img.Bounds())
    draw.Draw(rgba, img.Bounds(), img, image.Point{}, draw.Src)
    return rgba
}

// GlitchImage applies various glitch effects to the image
func (ig *ImageGlitcher) GlitchImage(glitchIntensity float32, cycle, step int, scanLines, pixelSort, colorOffset bool) (*image.RGBA, error) {
    rand.Seed(ig.Seed)
    for i := 0; i < cycle; i++ {
        if scanLines {
            applyScanLines(ig.Image, glitchIntensity)
        }
        if pixelSort {
            applyPixelSort(ig.Image)
        }
        if colorOffset {
            applyColorOffset(ig.Image, glitchIntensity)
        }
        applyRandomGlitch(ig.Image, glitchIntensity)
    }
    if isBlackImage(ig.Image) {
        return nil, fmt.Errorf(messages["blackImageError"])
    }
    return ig.Image, nil
}

// isBlackImage checks if the image is entirely black
func isBlackImage(img *image.RGBA) bool {
    for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
        for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
            r, g, b, a := img.At(x, y).RGBA()
            if r != 0 || g != 0 || b != 0 || a != 0 {
                return false
            }
        }
    }
    return true
}

// applyScanLines adds horizontal lines to the image
func applyScanLines(img *image.RGBA, glitchFactor float32) {
    step := int(max(1, glitchFactor))
    for y := 0; y < img.Bounds().Dy(); y += step {
        for x := 0; x < img.Bounds().Dx(); x++ {
            img.Set(x, y, color.RGBA{0, 0, 0, 255})
        }
    }
}

// applyPixelSort sorts the pixels in each row of the image based on their brightness
func applyPixelSort(img *image.RGBA) {
    for y := 0; y < img.Bounds().Dy(); y++ {
        row := make([]color.Color, img.Bounds().Dx())
        for x := 0; x < img.Bounds().Dx(); x++ {
            row[x] = img.At(x, y)
        }
        sort.Slice(row, func(i, j int) bool {
            ri, gi, bi, _ := row[i].RGBA()
            rj, gj, bj, _ := row[j].RGBA()
            return ri+gi+bi < rj+gj+bj
        })
        for x := 0; x < img.Bounds().Dx(); x++ {
            img.Set(x, y, row[x])
        }
    }
}

// applyColorOffset offsets the color channels of the image
func applyColorOffset(img *image.RGBA, glitchFactor float32) {
    width, height := img.Bounds().Dx(), img.Bounds().Dy()
    offsetX := rand.Intn(int(glitchFactor*2+1)) - int(glitchFactor)
    offsetY := rand.Intn(int(glitchFactor*2+1)) - int(glitchFactor)
    channel := rand.Intn(3)
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            newX := (x + offsetX + width) % width
            newY := (y + offsetY + height) % height
            r, g, b, a := img.RGBAAt(newX, newY).RGBA()
            switch channel {
            case 0:
                img.Set(x, y, color.RGBA{uint8(r >> 8), uint8(img.RGBAAt(x, y).G), uint8(img.RGBAAt(x, y).B), uint8(a >> 8)})
            case 1:
                img.Set(x, y, color.RGBA{uint8(img.RGBAAt(x, y).R), uint8(g >> 8), uint8(img.RGBAAt(x, y).B), uint8(a >> 8)})
            case 2:
                img.Set(x, y, color.RGBA{uint8(img.RGBAAt(x, y).R), uint8(img.RGBAAt(x, y).G), uint8(b >> 8), uint8(a >> 8)})
            }
        }
    }
}

// applyRandomGlitch applies random horizontal shifts to chunks of the image
func applyRandomGlitch(img *image.RGBA, glitchFactor float32) {
    maxOffset := int(glitchFactor * float32(img.Bounds().Dx()) / 10)
    for i := 0; i < int(glitchFactor*2); i++ {
        offset := rand.Intn(2*maxOffset) - maxOffset
        if offset == 0 {
            continue
        }
        if offset < 0 {
            glitchLeft(img, -offset)
        } else {
            glitchRight(img, offset)
        }
    }
}

// glitchLeft shifts a chunk of the image to the left
func glitchLeft(img *image.RGBA, offset int) {
    height := img.Bounds().Dy()
    startY := rand.Intn(height)
    chunkHeight := rand.Intn(height/4) + 1
    chunkHeight = min(chunkHeight, height-startY)
    stopY := startY + chunkHeight
    startX, stopX := offset, img.Bounds().Dx()-offset
    leftChunk := extractChunk(img, startY, stopY, startX, img.Bounds().Dx())
    wrapChunk := extractChunk(img, startY, stopY, 0, startX)
    placeChunk(img, leftChunk, startY, stopY, 0, stopX)
    placeChunk(img, wrapChunk, startY, stopY, stopX, img.Bounds().Dx())
}

// glitchRight shifts a chunk of the image to the right
func glitchRight(img *image.RGBA, offset int) {
    height := img.Bounds().Dy()
    startY := rand.Intn(height)
    chunkHeight := rand.Intn(height/4) + 1
    chunkHeight = min(chunkHeight, height-startY)
    stopY := startY + chunkHeight
    stopX, startX := img.Bounds().Dx()-offset, offset
    rightChunk := extractChunk(img, startY, stopY, 0, stopX)
    wrapChunk := extractChunk(img, startY, stopY, stopX, img.Bounds().Dx())
    placeChunk(img, rightChunk, startY, stopY, startX, img.Bounds().Dx())
    placeChunk(img, wrapChunk, startY, stopY, 0, startX)
}

// extractChunk extracts a chunk of the image
func extractChunk(img *image.RGBA, startY, stopY, startX, stopX int) *image.RGBA {
    chunk := image.NewRGBA(image.Rect(0, 0, stopX-startX, stopY-startY))
    for y := startY; y < stopY; y++ {
        for x := startX; x < stopX; x++ {
            chunk.Set(x-startX, y-startY, img.At(x, y))
        }
    }
    return chunk
}

// placeChunk places a chunk of the image
func placeChunk(img *image.RGBA, chunk *image.RGBA, startY, stopY, startX, stopX int) {
    for y := startY; y < stopY; y++ {
        for x := startX; x < stopX; x++ {
            img.Set(x, y, chunk.At(x-startX, y-startY))
        }
    }
}

// min returns the minimum of two integers
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// max returns the maximum of two float32 values
func max(a, b float32) float32 {
    if a > b {
        return a
    }
    return b
}

// saveImage saves the output image to the specified file, handling existing files
func saveImage(img image.Image, outputPath string) error {
    outFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf(messages["saveImageError"], err)
    }
    defer outFile.Close()

    switch strings.ToLower(filepath.Ext(outputPath)) {
    case ".png":
        err = png.Encode(outFile, img)
    case ".jpg", ".jpeg":
        err = jpeg.Encode(outFile, img, nil)
    default:
        return fmt.Errorf(messages["unsupportedFormat"], filepath.Ext(outputPath))
    }
    if err != nil {
        return fmt.Errorf(messages["encodingError"], err)
    }
    return nil
}

// saveGIF saves a series of images as a GIF
func saveGIF(images []*image.Paletted, delays []int, outputPath string) error {
    outFile, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf(messages["saveImageError"], err)
    }
    defer outFile.Close()

    g := &gif.GIF{
        Image: images,
        Delay: delays,
    }

    err = gif.EncodeAll(outFile, g)
    if err != nil {
        return fmt.Errorf(messages["encodingError"], err)
    }
    return nil
}

// spinner prints a cute flower spinner while a task is running
func spinner(done chan bool) {
    flowers := []string{"🌸", "🌺", "🌼", "🌷", "💐", "🌻", "🌹"}
    for {
        select {
        case <-done:
            fmt.Print(reset)
            return
        default:
            for _, f := range flowers {
                fmt.Printf("\r%s%s", pink, f)
                time.Sleep(200 * time.Millisecond)
            }
        }
    }
}

func generateRandomFilename(extension string) string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("glitched-image-%05d%s", rand.Intn(100000), extension)
}

func main() {
    // Define command-line flags
    glitchIntensity := flag.Float64("glitch-intensity", 5.0, "Intensity of the glitch effect (0.1-9.0)")
    scanLines := flag.Bool("scan-lines", false, "Apply scan lines glitch effect")
    pixelSort := flag.Bool("pixel-sort", false, "Apply pixel sort glitch effect")
    colorOffset := flag.Bool("color-offset", false, "Apply color offset glitch effect")
    seed := flag.Int("seed", -1, "Random seed for reproducibility")
    createGIF := flag.Bool("gif", false, "Create a GIF instead of a single image")
    frames := flag.Int("frames", 10, "Number of frames for the GIF")
    delay := flag.Int("delay", 10, "Delay between frames in the GIF")
    cycle := flag.Int("cycle", 1, "Number of cycles of glitches to apply")
    step := flag.Int("step", 1, "Step size between frames")

    // Parse flags
    flag.Parse()

    // Validate the number of arguments
    if len(flag.Args()) < 1 {
        fmt.Println(messages["usage"])
        os.Exit(0) // Exit with code 0 when showing help
    }

    // Read input and output file paths from arguments
    imgPath := flag.Arg(0)
    var outputPath string
    if len(flag.Args()) == 2 {
        outputPath = flag.Arg(1)
    } else {
        ext := filepath.Ext(imgPath)
        if *createGIF {
            outputPath = generateRandomFilename(".gif")
        } else {
            outputPath = generateRandomFilename(ext)
        }
    }

    // Validate input file
    if imgPath == "" {
        fmt.Println(messages["inputRequired"])
        os.Exit(0) // Exit with code 0 when showing help
    }

    // Validate glitchIntensity range
    if *glitchIntensity < 0.1 || *glitchIntensity > 9.0 {
        log.Fatalf("glitch-intensity must be between 0.1 and 9.0")
    }

    // Generate a random seed if none is provided
    randomSeed := time.Now().UnixNano()
    if *seed != -1 {
        randomSeed = int64(*seed)
    }

    // Initialize ImageGlitcher inline
    glitcher := &ImageGlitcher{
        Seed: randomSeed,
    }

    // Start spinner
    done := make(chan bool, 1)
    go spinner(done)
    fmt.Println(messages["startSpinner"])

    err := glitcher.FetchImage(imgPath)
    if err != nil {
        done <- true
        log.Fatalf(messages["fetchImageError"], err)
    }

    if *createGIF {
        // Create a GIF
        var images []*image.Paletted
        var delays []int
        for i := 0; i < *frames; i += *step {
            glitchedImg, err := glitcher.GlitchImage(float32(*glitchIntensity), *cycle, *step, *scanLines, *pixelSort, *colorOffset)
            if err != nil {
                done <- true
                log.Fatalf("%v\n", err)
            }
            palettedImg := image.NewPaletted(glitchedImg.Bounds(), palette.Plan9)
            draw.FloydSteinberg.Draw(palettedImg, glitchedImg.Bounds(), glitchedImg, image.Point{})
            images = append(images, palettedImg)
            delays = append(delays, *delay)
        }
        err = saveGIF(images, delays, outputPath)
        if err != nil {
            done <- true
            log.Fatalf("%v\n", err)
        }
    } else {
        // Create a single glitched image
        glitchedImg, err := glitcher.GlitchImage(float32(*glitchIntensity), *cycle, *step, *scanLines, *pixelSort, *colorOffset)
        if err != nil {
            done <- true
            log.Fatalf(messages["encodingError"], err)
        }
        if err := saveImage(glitchedImg, outputPath); err != nil {
            done <- true
            log.Fatalf(messages["saveImageError"], err)
        }
    }

    // Stop spinner
    done <- true
    fmt.Println(messages["stopSpinner"])
    fmt.Printf("%s%s%s\n", lightPink, messages["glitchComplete"], reset)
}

# glitcher
Glitcher is a re-implementation of the python image glitching tool `TotallyNotChase/glitch-this` which is based off the popular `image-glitcher` tool. This currently re-implements every feature of `glitch-this` including the creation of `*.gif` files. 

This is the starting point to create both the ability create videos from images with changes in glitching over time creating an animation, or taking every frame of a video and re-assembling the video with the glitched frames.

## Usage
Using the tool is very simple, as with any UI I develop, it is more fail-proof than the python version, you can provide as little as the input image path and it will fill in the rest for you, and if the file exists it will ask if you would like to overwrite it, and if you don't specify an output image path it will just use a randomly generated name prefix by `glitched-image-` and saved in the current working directory. 

```
Usage: glitcher [flags] input_image [output_image]
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
```


# CheckBoxDetector
Software for detect checkboxes in form images

## Context
The task of detect filled checkboxes in an image was a really challenging for me since I have 
no experience working with images.

## Architecture
The code follow hexagonal architecture with the definition of ports & corresponding adapters.

I think it this way to make it scalable maybe with integration of other APIs that solve some of the faced problems 
(get images, decoding it or apply thresholds to binarize it)

## Briefly explanation of code
The code it's a sequential orchestration of steps that consist on:

- Get the image
- Decode the image
- Convert the image to gray scale
- Binarize the image 
- Detect contours
- Filter possible rectangles
- Check if the rectangle is filled or not

## How to run it

In the root of the project, run one of the following commands: ``go run cmd/main.go`` or ``make run``

## Output
The output of the development is the position of the detected checkboxes printed in the console and 
the number of total checkboxes detected

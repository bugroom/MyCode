package main

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())

	fmt.Println("Something you can do")
	fmt.Println("A: Create a random rect on a image")
	fmt.Println("B: Merge two pictures into one")

	var s string
	fmt.Scanln(&s)

	if s == "A" || s == "a" {
		Random_color_rect()
	} else if s == "B" || s == "b" {
		Merge()
	}

}

func decode_img(f *os.File, f_name string) image.Image {
	var img image.Image
	var err error

	if f_name[len(f_name)-3:] == "jpg" || f_name[len(f_name)-4:] == "jpeg" {
		img, err = jpeg.Decode(f)
		if err != nil {
			fmt.Println(err)
		}
	} else if f_name[len(f_name)-3:] == "png" {
		img, err = png.Decode(f)
		if err != nil {
			fmt.Println(err)
		}
	}
	return img
}

func encode_img(f *os.File, f_name string, img image.Image) {
	var err error
	if f_name[len(f_name)-3:] == "jpg" || f_name[len(f_name)-4:] == "jpeg" {
		err = jpeg.Encode(f, img, nil) //nil for DefaultQuality
		if err != nil {
			fmt.Println(err)
		}
	} else if f_name[len(f_name)-3:] == "png" {
		err = png.Encode(f, img)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Random_color_rect() {
	fmt.Println("Input the name of the original image (only accept .png/.jpg/.jpeg file)") //read a img file
	var file_name string
	fmt.Scanln(&file_name)
	src_file, err := os.Open(file_name)
	if err != nil {
		fmt.Println(err)
	}
	defer src_file.Close()
	img := decode_img(src_file, file_name) //turn the img file into data

	color_pic := image.NewRGBA(image.Rect(0, 0, 20+rand.Intn(img.Bounds().Dx()), 20+rand.Intn(img.Bounds().Dy())))
	//color size is random and smaller than src img, +20 is to avoid the pic being too small to see
	color := palette.Plan9[rand.Intn(256)]
	//random color
	draw.Draw(color_pic, color_pic.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	//fill color//image.ZP is Point{0,0}
	img_draw := image.NewRGBA(img.Bounds())
	draw.Draw(img_draw, img_draw.Bounds(), img, image.ZP, draw.Src)
	//turn image.Image into draw.Image, cuz draw.Draw()'s fist argument has to be draw.Image
	x0 := rand.Intn(img_draw.Bounds().Dx())
	y0 := rand.Intn(img_draw.Bounds().Dy())
	draw.Draw(img_draw, image.Rect(x0, y0, x0+color_pic.Bounds().Dx(), y0+color_pic.Bounds().Dy()), color_pic, image.ZP, draw.Src)
	//in a random place, draw the random color rect

	fmt.Println("Input the name of the new image (only accept .png/.jpg/.jpeg file)") //create the new file
	fmt.Scanln(&file_name)
	dst_file, err := os.Create(file_name)
	if err != nil {
		fmt.Println(err)
	}
	defer dst_file.Close()
	encode_img(dst_file, file_name, img_draw)
}

func Merge() {
	var file_name_1, file_name_2, file_result string

	fmt.Println("Input the name of the two original pictures, separate with space key (only accept .png/.jpg/.jpeg file)")
	fmt.Scanln(&file_name_1, &file_name_2)
	//open file and decode the two pics
	f1, err := os.Open(file_name_1)
	if err != nil {
		fmt.Println(err)
	}
	defer f1.Close()
	f2, err := os.Open(file_name_2)
	if err != nil {
		fmt.Println(err)
	}
	defer f2.Close()
	img1 := decode_img(f1, file_name_1)
	img2 := decode_img(f2, file_name_2)

	//find the smallest bound (the common part)of the two
	var x_min, y_min int
	if img1.Bounds().Dx() < img2.Bounds().Dx() {
		x_min = img1.Bounds().Dx()
	} else {
		x_min = img2.Bounds().Dx()
	}
	if img1.Bounds().Dy() < img2.Bounds().Dy() {
		y_min = img1.Bounds().Dy()
	} else {
		y_min = img2.Bounds().Dy()
	}

	dst := image.NewRGBA(image.Rect(0, 0, x_min, y_min)) //the bound of dst should be the common part

	//fill dst, every pixel is random to be from src1 or src2
	for x := 0; x < x_min; x++ {
		for y := 0; y < y_min; y++ {
			if rand.Intn(2) == 1 {
				dst.Set(x, y, img1.At(x, y))
			} else {
				dst.Set(x, y, img2.At(x, y))
			}
		}
	}

	fmt.Println("Input the name of the new image file (only accept .png/.jpg/.jpeg file)")
	fmt.Scanln(&file_result)
	f_dst, err := os.Create(file_result)
	if err != nil {
		fmt.Println(err)
	}
	defer f_dst.Close()
	encode_img(f_dst, file_result, dst)
}
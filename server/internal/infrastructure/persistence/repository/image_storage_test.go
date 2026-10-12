package repository_test

// func TestImageRepository(t *testing.T) {
// 	container := di.BuildContainer()
// 	ctx := context.Background()
// 	repo := di.Invoke[image.ImageRepository](container)
// 	dummyImage, _ := createSimpleDummyJPEG(1, 1)
// 	imageInfo := image.NewImageInfo("test/hoge.jpg", "jpg", image.NewImage(dummyImage))

// 	repo.Create(ctx, *imageInfo)
// }

// func createSimpleDummyJPEG(width, height int) ([]byte, error) {
// 	// 単色の画像を作成
// 	img := img.NewRGBA(img.Rect(0, 0, width, height))
// 	blue := color.RGBA{0, 0, 255, 255} // 青色
// 	for y := 0; y < height; y++ {
// 		for x := 0; x < width; x++ {
// 			img.Set(x, y, blue)
// 		}
// 	}

// 	// JPEGにエンコード
// 	var buf bytes.Buffer
// 	err := jpeg.Encode(&buf, img, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return buf.Bytes(), nil
// }

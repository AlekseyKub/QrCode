// QrCode
package main

import (
	"fmt"

	"github.com/AlekseyKub/qrgen"
)

func main() {

	var VersionCode int = 3 //Версия кода
	var correctionLevel = 1 //Уровень коррекциих[0-3] 0-L, 1-M, 2-Q, 3-H
	var codeIndicators = 4  //Индикатор режима
	var maskNumber = 0      //Номер маски

	//срез для данных окончательного заполнения QR кода
	var dataOut = make([]byte, 0, (qrgen.MaxDataBit[VersionCode-1][correctionLevel]/8)+
		(qrgen.NumberOfBlocks[VersionCode-1][correctionLevel]*
			qrgen.NumberOfBytesCorrection[VersionCode-1][correctionLevel]))

	//Входные данные
	var dataIn = []byte("https://google.com")
	//Создаю срез для обработки данных
	var data = make([]byte, 0, qrgen.MaxDataBit[VersionCode-1][correctionLevel]/8)

	//Проверка помещются ли данные в выбранный QR код
	if qrgen.SizeСheck(dataIn, VersionCode-1, correctionLevel) {
		data = append(data, byte(codeIndicators))
		data = append(data, byte(len(dataIn)))
		data = append(data, dataIn...)
		bitShift(data, 4) //сдвигаю данные на 4 бита влево
	} else {
		fmt.Println("Увеличить версию")
	}

	//Дополняю данные до нужной длинны байтами 236 17
	data = fillData(data, VersionCode-1, correctionLevel)

	//Формирую массив для заполнения QR
	//(расчет блоков коррекции + заполение данными + блокими коррекции)
	dataOut = qrgen.AddingDataOut(dataOut, data, VersionCode-1, correctionLevel)

	// Создание двумерного среза для хранения QR кода
	qrCode := make([][]int, qrgen.SizeQrCode[VersionCode-1])
	for i := range qrCode {
		qrCode[i] = make([]int, qrgen.SizeQrCode[VersionCode-1])
	}

	//Нанесение поисковых узоров
	qrgen.MakeSearchPattern(qrCode)
	//Нанесение выравнивающих узоров
	qrgen.CreatePattern(qrgen.LevelingPattern[VersionCode-1], qrCode)
	//Нанесение маски и уровня коррекции
	qrgen.CreateMaskAndCorrectionLevel(qrCode, VersionCode, correctionLevel, maskNumber)

	//Заполнение QR кода данными
	qrgen.AddingDataQR(qrCode, dataOut, VersionCode-1)
	//Применение маски
	qrgen.ApplyMask(qrCode, maskNumber)

	//Вывести Qr код

	// for i := 0; i < (len(qrCode)); i++ {
	// 	for x := 0; x < (len(qrCode)); x++ {
	// 		fmt.Print(qrCode[i][x])
	// 		if x == (len(qrCode))-1 {
	// 			fmt.Println()
	// 		}
	// 	}
	// }

	//создать SVG файл
	qrgen.CreateSvg(qrCode)
}

func printBit(x int) {
	for i := 7; i >= 0; i-- {
		if hasBit(x, i) {
			fmt.Print(1)
		} else {
			fmt.Print(0)
		}
	}
	fmt.Print(" ")
}

// Проверить бит (если 1 вертнет TRUE)
func hasBit(n int, pos int) bool {
	val := n & (1 << pos)
	return (val > 0)
}

// Сдвиг байтов массива данных
func bitShift(arr []byte, x int) {
	for i := 0; i < len(arr); i++ {
		if i < len(arr)-1 {
			arr[i] = (arr[i] << x) | (arr[i+1] >> x)
		} else {
			arr[i] = arr[i] << x
		}
	}
}

// Дополнить данный до нужной длинны
func fillData(data []byte, ver int, cor int) []byte {
	a := 236 //11101100
	b := 17  //00010001
	x := a
	maxData := qrgen.MaxDataBit[ver][cor] / 8
	for i := maxData - len(data); i != 0; i-- {
		data = append(data, byte(x))
		if x == a {
			x = b
		} else {
			x = a
		}
	}
	return data
}

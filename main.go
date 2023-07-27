// QrCode
package main

import (
	"fmt"
	"os"

	"github.com/AlekseyKub/qrgen"
)

var VersionCode int = 0     //Версия кода
var correctionLevel int = 0 //Уровень коррекциих[0-3] 0-L, 1-M, 2-Q, 3-H
var codeIndicators int = 4  //Индикатор режима
var maskNumber int = 0      //Номер маски
var fileName string = ""    //Название файла

func main() {

	//Входные данные
	var dataIn = []byte("")

	//argsWithoutProg := os.Args[1:]
	dataIn = make([]byte, 0, len(os.Args[1]))
	dataIn = append(dataIn, os.Args[1]...)
	fileName = os.Args[2]

	arg3 := os.Args[3]
	switch arg3 {
	case "l":
		correctionLevel = 0
	case "m":
		correctionLevel = 1
	case "q":
		correctionLevel = 2
	case "h":
		correctionLevel = 3
	}

	//Определение необходимой версии
	VersionCode = qrgen.CalculatVersion(len(dataIn), correctionLevel)
	fmt.Println(VersionCode)

	//срез для данных окончательного заполнения QR кода
	var dataOut = make([]byte, 0, (qrgen.MaxDataBit[VersionCode-1][correctionLevel]/8)+
		(qrgen.NumberOfBlocks[VersionCode-1][correctionLevel]*
			qrgen.NumberOfBytesCorrection[VersionCode-1][correctionLevel]))

	//Создаю срез для обработки данных
	var data = make([]byte, 0, qrgen.MaxDataBit[VersionCode-1][correctionLevel]/8)

	//Проверка помещются ли данные в выбранный QR код
	if qrgen.SizeСheck(dataIn, VersionCode-1, correctionLevel) {
		data = append(data, byte(codeIndicators))
		if VersionCode < 10 { //если меньше 10 версии
			data = append(data, byte(len(dataIn))) //поле количество данных 1 байт
		} else {
			data = append(data, (howMuchData(len(dataIn)))...) //поле количества данных 2 байта
		}
		data = append(data, dataIn...)
		bitShift(data, 4) //сдвигаю данные на 4 бита влево
	} else {
		fmt.Println("Увеличить версию QR кода")
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
	if VersionCode > 1 {
		qrgen.CreatePattern(qrgen.LevelingPattern[VersionCode-1], qrCode)
	}

	//Нанесение шаблона синхронизации (верт. и гор. зонатьной послендовательности 1 и 0)
	qrgen.СreateTimingTemplate(qrCode)

	//Нанесение маски и уровня коррекции
	qrgen.CreateMaskAndCorrectionLevel(qrCode, VersionCode-1, correctionLevel, maskNumber)

	//Заполнение QR кода данными
	qrgen.AddingDataQR(qrCode, dataOut, VersionCode-1)

	//Применение маски
	qrgen.ApplyMask(qrCode, maskNumber)

	// //Вывести Qr код
	// for i := 0; i < (len(qrCode)); i++ {
	// 	for x := 0; x < (len(qrCode)); x++ {
	// 		fmt.Print(qrCode[i][x])
	// 		if x == (len(qrCode))-1 {
	// 			fmt.Println()
	// 		}
	// 	}
	// }

	//создать SVG файл
	qrgen.CreateSvg(qrCode, fileName)
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

// Дополнить данные до нужной длинны
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

// раделение на 2 байта количсво данных (для версии > 10)
func howMuchData(num int) []byte {
	quantity := []byte{0, 0}
	quantity[0] = byte(num >> 8)
	quantity[1] = byte(num)

	return quantity
}

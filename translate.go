package main

func translateDataTypeFS(input string) string {
	var output string
	switch input {
	case "0":
		output = "0"
	case "1":
		output = "1"
	case "2":
		output = "5"
	case "3":
		output = "2"
	case "4":
		output = "6"
	case "5":
		output = "3"
	case "6":
		output = "7"
	case "7":
		output = "1"
	case "8":
		output = "3"
	case "9":
		output = "5"
	case "10":
		output = "9"
	case "11":
		output = "TEXT"
	case "12":
		output = "4"
	case "13":
		output = "8"
	case "14":
		output = "10"
	default:
		output = "NaN"
	}

	return output
}

func translateDataTypeAS(dataTypeFS string, DdataTypeAS string) string {
	var output string
	switch dataTypeFS {
	case "7":
		output = "2"
	case "8":
		output = "2"
	case "9":
		output = "2"
	default:
		output = DdataTypeAS
	}

	return output
}

func loytecToModulo5DataTypeAS(registerType string) string {
	var output string
	switch registerType {
	case "INPUT":
		output = "2"
	case "DISCRETE":
		output = "2"
	case "HOLD":
		output = "0"
	default:
		output = "0"
	}

	return output
}

func loytecToModulo5DataTypeFS(dataType string) string {
	var output string
	switch dataType {
	case "bit":
		output = "0"
	case "uint8":
		output = "1"
	case "uint16":
		output = "3"
	case "uint32":
		output = "5"
	case "uint64":
		output = "12"
	case "int8":
		output = "2"
	case "int16":
		output = "4"
	case "int32":
		output = "6"
	case "int64":
		output = "13"
	case "float32":
		output = "10"
	case "float64":
		output = "14"
	default:
		output = "3"
	}

	return output
}

func modulo5ToLoytecRegisterLength(dataTypeFS string) string {
	var output string
	switch dataTypeFS {
	case "0":
		output = "0"
	case "1":
		output = "1"
	case "3":
		output = "2"
	case "5":
		output = "4"
	case "12":
		output = "8"
	case "2":
		output = "1"
	case "4":
		output = "2"
	case "6":
		output = "4"
	case "13":
		output = "8"
	case "10":
		output = "4"
	case "14":
		output = "8"
	default:
		output = "1"
	}

	return output
}

func modulo5ToLoytecModbusDataType(dataTypeFS string) string {
	var output string
	switch dataTypeFS {
	case "0":
		output = "bit"
	case "1":
		output = "uint8"
	case "3":
		output = "uint16"
	case "5":
		output = "uint32"
	case "12":
		output = "uint64"
	case "2":
		output = "int8"
	case "4":
		output = "int16"
	case "6":
		output = "int32"
	case "13":
		output = "in64"
	case "10":
		output = "float32"
	case "14":
		output = "float64"
	default:
		output = "uint16"
	}

	return output
}

func loytecToModulo5FunctionCode(registerType string, Direction string) string {
	output := "0"

	if registerType == "DISCRETE" {
		output = "2"
	}

	if registerType == "COIL" {
		if Direction == "Input" {
			output = "1"
		} else if Direction == "Output" {
			output = "5"
		} else if Direction == "Value" {
			output = "5"
		}
	}

	if registerType == "INPUT" {
		output = "4"
	}

	if registerType == "HOLD" {
		if Direction == "Input" {
			output = "3"
		} else if Direction == "Output" {
			output = "6"
		} else if Direction == "Value" {
			output = "6"
		}
	}

	return output
}

func modulo5ToLoytecRegisterType(functionCode string) string {
	var output string
	switch functionCode {
	case "2":
		output = "DISCRETE"
	case "1":
		output = "COIL"
	case "4":
		output = "INPUT"
	case "3":
		output = "HOLD"
	case "5":
		output = "COIL"
	case "6":
		output = "HOLD"
	default:
		output = "INPUT"
	}

	return output
}

// Just a mock translation for future implementation
func translateByteOrder(byteOrder string) string {
	return "1"
}

// Just a mock translation for future implementation
func translateWordOrder(byteOrder string) string {
	return byteOrder
}

func loytecToModulo5ByteOrder(swap16, swap32, swap64 string) string {
	if swap16 == "1" && swap32 == "0" && swap64 == "0" {
		return "0"
	}
	if swap16 == "0" && swap32 == "1" && swap64 == "0" {
		return "1"
	}
	if swap16 == "0" && swap32 == "0" && swap64 == "1" {
		return "2"
	}
	return "0"
}

func modulo5ToLoytecSwap(byteOrder string) (string, string, string) {
	if byteOrder == "0" {
		return "1", "0", "0"
	}
	if byteOrder == "1" {
		return "0", "1", "0"
	}
	if byteOrder == "2" {
		return "0", "0", "1"
	}
	return "0", "0", "0"
}

func loytecToModulo5ScalingA(exponent string) string {
	var output string
	switch exponent {
	case "-3":
		output = "1000"
	case "-2":
		output = "100"
	case "-1":
		output = "10"
	case "0":
		output = "1"
	case "1":
		output = "0.1"
	case "2":
		output = "0.01"
	case "3":
		output = "0.001"
	default:
		output = "1"
	}

	return output
}

func modulo5ToLoytecExponent(ScalingA string) string {
	var output string
	switch ScalingA {
	case "1000":
		output = "-3"
	case "100":
		output = "-2"
	case "10":
		output = "-1"
	case "1":
		output = "0"
	case "0.1":
		output = "1"
	case "0.01":
		output = "2"
	case "0.001":
		output = "3"
	default:
		output = "0"
	}

	return output
}

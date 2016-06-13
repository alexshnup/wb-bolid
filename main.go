package main

import (
        "log"
        "fmt"
        "github.com/tarm/serial"
        "time"
)

// CRC 8-bit Calculate method Dallas 1-wire with prepared tables
func crc8dallas(Command1 []byte) []byte{
  CRCTable := [256]byte{0,94,188,226,97,63,221,131,194,156,126,32,163,253,31,65,157,195,33,127,252,162,64,30,95,1,227,189,62,96,130,220,35,125,159,193,66,28,254,160,225,191,93,3,128,222,60,98,190,224,2,92,223,129,99,61,124,34,192,158,29,67,161,255,70,24,250,164,39,121,155,197,132,218,56,102,229,187,89,7,219,133,103,57,186,228,6,88,25,71,165,251,120,38,196,154,101,59,217,135,4,90,184,230,167,249,27,69,198,152,122,36,248,166,68,26,153,199,37,123,58,100,134,216,91,5,231,185,140,210,48,110,237,179,81,15,78,16,242,172,47,113,147,205,17,79,173,243,112,46,204,146,211,141,111,49,178,236,14,80,175,241,19,77,206,144,114,44,109,51,209,143,12,82,176,238,50,108,142,208,83,13,239,177,240,174,76,18,145,207,45,115,202,148,118,40,171,245,23,73,8,86,180,234,105,55,213,139,87,9,235,181,54,104,138,212,149,203,41,119,244,170,72,22,233,183,85,11,136,214,52,106,43,117,151,201,74,20,246,168,116,42,200,150,21,75,169,247,182,232,10,84,215,137,107,53}
  size := len(Command1) - 1
  Command1[size] = 0
  for i := 0; i <= int(size) - 1; i++ {
    Command1[size] = CRCTable[Command1[i] ^ Command1[size]]
  }
return Command1
}

//Send and Receive Answer bytes
func write_serial(send []byte) []byte{
  delay := 50 * time.Millisecond

  //c := &serial.Config{Name: "/dev/ttyS2", Baud: 9600, ReadTimeout: time.Millisecond * 500}
  c := &serial.Config{Name: "/dev/ttyAPP4", Baud: 9600, ReadTimeout: time.Millisecond * 5000}

   //c := new(serial.Config)
   //c.Name = "/dev/ttyAPP1"
   //c.Name = "/dev/ttyS2"
   //c.Baud = 9600
   //c.ReadTimeout = time.Millisecond * 5000
   //c.Size = 8
   //c.StopBits = 0

  s, err := serial.OpenPort(c)
  if err != nil {
          log.Fatal(err)
  }
  //Send Data
  //n, err := s.Write([]byte("test"))
  //s.Flush()
  time.Sleep(delay)
  _, err = s.Write(crc8dallas(send))
  if err != nil {
	  log.Printf("Error Oper Port...")
          log.Fatal(err)
  }

  //Need delay for correct the receive answer
  time.Sleep(delay)

  //Receive Respond
  buf := make([]byte, 6)
  _, err = s.Read(buf)
  //nr, err_read := s.Read(buf)
  if err != nil {
	  log.Printf("Error Read...")
          log.Fatal(err)
  }

  //Discards data written to the port but not transmitted,
  //or data received but not read
  s.Flush()

  //Close Serial Port
  s.Close()

  return []byte(buf)
}

func status_serial(result []byte, ok uint8, relay uint8) {
  fmt.Printf("%x\n", result)
  switch ok {
  case 0:
    if (rune(result[2]) == 0x18) && (rune(result[5]) == 0xEC){
      log.Printf("Enter - relay %d - OK", relay)
    }
  case 1:
    if (rune(result[2]) == 0x42) && (rune(result[5]) == 0xFD){
      log.Printf("Command OK - relay %d", relay)
    }
  case 2:
    if (rune(result[4]) == 0x01) && (rune(result[5]) == 0xB6){
      log.Printf("Port %d - ON", relay)
    } else if (rune(result[4]) == 0x02) && (rune(result[5]) == 0x54) {
      log.Printf("Port %d - OFF", relay)
    }


  }
}


func Relay_ON (addr uint8, relay uint8) {
  Command1 := []byte{127,8,0,65,1,0,0,1,0}
  CommandEnter := []byte{127,6,0,23,0,0,0}
  Command1[0] = addr
  Command1[4] = relay
  CommandEnter[0] = addr
  status_serial(write_serial(crc8dallas(Command1)), 1, relay)
  status_serial(write_serial(crc8dallas(CommandEnter)), 0, relay)
}

func Relay_OFF (addr uint8, relay uint8) {
  Command1 := []byte{127,8,0,65,1,0,0,2,0}
  CommandEnter := []byte{127,6,0,23,0,0,0}
  Command1[0] = addr
  Command1[4] = relay
  CommandEnter[0] = addr
  status_serial(write_serial(crc8dallas(Command1)), 1, relay)
  status_serial(write_serial(crc8dallas(CommandEnter)), 0, relay)
}

func StatusRelay (addr uint8, relay uint8) {
  Command1 := []byte{127,8,0,67,1,0,0,1,0}
  Command1[0] = addr
  Command1[4] = relay
  status_serial(write_serial(crc8dallas(Command1)), 2, relay)
}

func main() {

      Relay_ON(127, 1)
      StatusRelay(127, 1)
        time.Sleep(100 * time.Millisecond)
      Relay_OFF(127, 1)
      StatusRelay(127, 1)

      Relay_ON(127, 2)
      StatusRelay(127, 2)
        time.Sleep(100 * time.Millisecond)
      Relay_OFF(127, 2)
      StatusRelay(127, 2)

      Relay_ON(127, 3)
      StatusRelay(127, 3)
        time.Sleep(100 * time.Millisecond)
      Relay_OFF(127, 3)
      StatusRelay(127, 3)

      Relay_ON(127, 4)
      StatusRelay(127, 4)
        time.Sleep(100 * time.Millisecond)
      Relay_OFF(127, 4)
      StatusRelay(127, 4)

}


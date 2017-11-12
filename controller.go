package xpider

import (
	"github.com/celer/xpider/hdlc"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	HEARTBEAT       uint8 = 0x00
	MOVE            uint8 = 0x01
	FRONT_LED       uint8 = 0x02
	EYE             uint8 = 0x03
	STOP            uint8 = 0x08
	GET_REG         uint8 = 0x09
	UPDATE_REG      uint8 = 0x0A
	REGISTER_RETURN uint8 = 0x0B
	WALK_BY_STEP    uint8 = 0x0D
	AUTO_MOVE       uint8 = 0x0E
)

// RobotState holds the last state of the robot
type RobotState struct {
	// StepCounter (0 - 65535)
	StepCounter uint16
	// ObsticalDistance (0 - 500) (500 if invalid)
	ObsticalDistance uint16
	// BatteryVoltage is the voltage in 10mV(1 - 10000)
	BatteryVoltage uint16
	// Yaw ( 0 - 2PI )
	Yaw float32
	// Pitch ( 0 - 2PI )
	Pitch float32
	// Roll ( 0 - 2PI )
	Roll float32
	// Sound ( 0 - 255 )
	Sound uint8
	// Updated when the robot's state was last updated
	Updated time.Time
}

// Controller controls the robot
type Controller struct {
	state       RobotState
	UpdateMutex sync.Mutex
	Conn        net.Conn
}

// Connect to the robot
//	addr - IP:Port address
func (x *Controller) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	x.Conn = conn

	go func() {
		k := make([]byte, 2048)
		b := make([]byte, 2048)
		for true {
			n, err := x.Conn.Read(b)
			if err != nil {
				panic(err)
			}
			if b[0] == 0x01 && b[1] == 0x55 {
				r := hdlc.NewReader(bytes.NewBuffer(b[2:]))
				s, err := r.Read(k)
				if err != nil {
					fmt.Printf("Error reading data %v\n", err)
				}
				if k[0] == HEARTBEAT {
					x.heartBeat(k[:s])
				}
			} else {
				fmt.Printf("?? %x\n", b[:n])
			}
		}
	}()

	return nil
}

// GetState gets the last reported state of the robot
func (x *Controller) GetState() RobotState {
	x.UpdateMutex.Lock()
	defer x.UpdateMutex.Unlock()
	return x.state
}

func (x *Controller) heartBeat(k []byte) {
	x.UpdateMutex.Lock()
	defer x.UpdateMutex.Unlock()

	x.state.StepCounter = binary.LittleEndian.Uint16(k[1:])
	x.state.ObsticalDistance = binary.LittleEndian.Uint16(k[3:])
	x.state.BatteryVoltage = binary.LittleEndian.Uint16(k[5:])

	x.state.Sound = k[19]

	gyro := bytes.NewReader(k[7:])
	binary.Read(gyro, binary.LittleEndian, &x.state.Yaw)
	binary.Read(gyro, binary.LittleEndian, &x.state.Pitch)
	binary.Read(gyro, binary.LittleEndian, &x.state.Roll)

	x.state.Updated = time.Now()
}

// Write a raw chunk of data to the robot
func (x *Controller) Write(data []byte) (int, error) {
	b := &bytes.Buffer{}
	b.Write([]byte{0x01, 0x55})
	w := hdlc.NewWriter(b)
	w.Write(data)
	return x.Conn.Write(b.Bytes())
}

// WalkByStep
//	speed (-100 - 100)
//	stepNum (2 - 255)
func (x *Controller) WalkByStep(speed int8, stepNum uint8) error {
	_, err := x.Write([]byte{WALK_BY_STEP, byte(speed), byte(stepNum)})
	return err
}

// Move the robot
// 	walkSpeed (-100 - 100)
// 	rotate (-100 - 100)
func (x *Controller) Move(walkSpeed int8, rotate int8) error {
	_, err := x.Write([]byte{MOVE, byte(walkSpeed), byte(rotate)})
	return err
}

// FrontLED powers on the front LEDS
// 	Left RGB value 0 - 225 for each component
//	Right RGB value 0 - 255 for each component
func (x *Controller) FrontLED(lr, lg, lb, rr, rg, rb uint8) error {
	_, err := x.Write([]byte{FRONT_LED, lr, lg, lb, rr, rg, rb})
	return err
}

// Eye rotates the eye up and down
// 	cameraEnabled (disabled)
// 	eyeAngle (15 - 65)
func (x *Controller) Eye(cameraEnable byte, eyeAngle uint8) error {
	_, err := x.Write([]byte{EYE, cameraEnable, eyeAngle})
	return err

}

// AutoMove moves the robot
//	rotateSpeed 	(0 - 100)
// 	rotateAngle 	(-PI - PI)
//	walkSpeed	(-100 - 100)
// 	stepNum		(2 - 255)
func (x *Controller) AutoMove(rotateSpeed uint8, rotateAngle float32, walkSpeed int8, stepNum uint8) error {
	b := &bytes.Buffer{}
	b.Write([]byte{AUTO_MOVE, rotateSpeed})
	binary.Write(b, binary.LittleEndian, rotateAngle)
	b.Write([]byte{byte(walkSpeed)})
	b.Write([]byte{byte(stepNum)})
	_, err := x.Write(b.Bytes())
	return err
}

# Thorlabs KDC101

A Go library for communicating with Thorlabs KDC101 DC Servo Motor Controllers. This library provides a comprehensive interface for controlling precision motion stages including translation stages, rotation mounts, and cage rotation systems.

## Features

- **Multi-Protocol Support**: Communicate via Serial (RS-232/USB) using the unified Unicomm interface
- **Precision Motion Control**: Full support for absolute positioning, relative moves, continuous motion, and jogging
- **Stage-Specific Calibration**: Automatic unit conversions for supported Thorlabs stages and motor types
- **Real-Time Status**: Monitor position, velocity, current, and comprehensive status flags
- **Velocity Profiling**: Configure trapezoidal velocity profiles for smooth motion control
- **Hardware Identification**: Query controller information and enable LED identification

## Installation

```bash
go get github.com/devicehub-go/thorlabs-kdc101
```

## Supported Hardware

### Motor Controllers
- **KDC101** - 1-Channel DC Servo Motor Controller

### Motor Types
- Brushed DC motors
- Brushless DC motors

### Supported Stages
- **MTS25-Z8** - 25mm Translation Stage
- **MTS50-Z8** - 50mm Translation Stage  
- **Z8xx Series** - Various Z8 stages
- **Z6xx Series** - Various Z6 stages
- **PRM1-Z8** - Rotation Mount
- **PRMTZ8** - Rotation Mount
- **CR1-Z7** - Cage Rotation Mount
- **KVS30** - Vertical Translation Stage

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    kdc101 "github.com/devicehub-go/thorlabs-kdc101"
    "github.com/devicehub-go/unicomm"
)

func main() {
    // Create controller instance
    controller := kdc101.New(
        kdc101.MTS25Z8,    // Stage type
        kdc101.Brushed,    // Motor ty
            Port: "/dev/ttyUSB0",  // Adjust for your system
	        BaudRate: 115200,
			DataBits: 8,
			StopBits: 1,
			Parity: "N",
            ReadTimeout:  time.Second * 5,
            WriteTimeout: time.Second * 5,
        },
    )
    
    // Connect to the controller
    if err := controller.Connect(); err != nil {
        log.Fatal("Failed to connect:", err)
    }
    defer controller.Disconnect()
    
    // Enable the motor
    if err := controller.Enable(1, true); err != nil {
        log.Fatal("Failed to enable motor:", err)
    }
    
    // Home the stage
    fmt.Println("Homing stage...")
    if err := controller.StartHomeMove(1); err != nil {
        log.Fatal("Failed to start homing:", err)
    }
    
    // Wait for homing to complete
    for {
        status, err := controller.GetDCStatusUpdate(1)
        if err != nil {
            log.Fatal("Failed to get status:", err)
        }
        
        statusSI := controller.DCStatusUpdateToSI(status)
        if statusSI.StatusBits.IsHomed {
            fmt.Println("Homing complete!")
            break
        }
        
        time.Sleep(100 * time.Millisecond)
    }
    
    // Move to 10mm position
    fmt.Println("Moving to 10mm...")
    if err := controller.MoveAbsolutePosition(1, 10.0); err != nil {
        log.Fatal("Failed to move:", err)
    }
}
```

## API Reference

### Constructor

#### `New(stage StageType, motor MotorType, options unicomm.UnicommOptions) *KDC101`
Creates a new KDC101 instance with the specified stage type, motor type, and communication options.

**Parameters:**
- `stage`: Stage type constant (e.g., MTS25Z8, MTS50Z8)
- `motor`: Motor type constant (Brushed or Brushless)
- `options`: Communication configuration

### Connection Management

#### `Connect() error`
Establishes connection with the controller and initializes communication.

#### `Disconnect() error`
Closes the connection with the controller.

### Device Information

#### `GetInformation() (HwInformation, error)`
Returns comprehensive hardware information including serial number, model, firmware version, and channel count.

#### `Identify(channel uint8) error`
Instructs the controller to flash its front panel LEDs for identification. Channel must be 1.

### Motor Control

#### `Enable(channel uint8, enable bool) error`
Enables or disables the specified motor channel.

#### `IsEnabled(channel uint8) (bool, error)`
Returns the enabled state of the motor channel.

#### `Stop(channel uint8, mode StopMode) error`
Stops motor motion using the specified stop mode (Abrupt or Soft).

### Motion Commands

#### `StartHomeMove(channel uint8) error`
Initiates a homing sequence using the configured home parameters.

#### `MoveAbsolutePosition(channel uint8, position float64) error`
Moves to an absolute position in millimeters.

#### `StartAbsoluteMove(channel uint8) error`
Starts an absolute move using previously set parameters.

#### `MoveRelativeDistance(channel uint8, distance float64) error`
Moves a relative distance in millimeters from the current position.

#### `StartRelativeMove(channel uint8) error`
Starts a relative move using previously set parameters.

#### `StartJogMove(channel uint8, direction Direction) error`
Performs a jog move in the specified direction (Forward or Reverse).

#### `MoveContinuous(channel uint8, direction Direction) error`
Moves continuously in the specified direction until stopped or a limit is reached.

### Parameter Configuration

#### `SetTrapezoidalVelocity(channel uint8, profile VelocityProfile) error`
Sets trapezoidal velocity parameters for smooth motion control.

```go
profile := VelocityProfile{
    MinVelocity:  0.1,    // mm/s
    MaxVelocity:  5.0,    // mm/s
    Acceleration: 10.0,   // mm/s²
}
```

#### `GetTrapezoidalVelocity(channel uint8) (VelocityProfile, error)`
Returns the current trapezoidal velocity parameters.

#### `SetJogParameters(channel uint8, params JogParameters) error`
Configures jog motion parameters including step size, velocities, and acceleration.

#### `GetJogParameters(channel uint8) (JogParameters, error)`
Returns the current jog parameters.

#### `SetRelativeMoveDistance(channel uint8, distance float64) error`
Sets the distance for the next relative move operation.

#### `GetRelativeMoveDistance(channel uint8) (float64, error)`
Returns the configured relative move distance.

#### `SetAbsoluteMoveDistance(channel uint8, position float64) error`
Sets the target position for the next absolute move operation.

#### `GetAbsoluteMoveDistance(channel uint8) (float64, error)`
Returns the configured absolute move target position.

### Status Monitoring

#### `GetDCStatusUpdate(channel uint8) (DCStatusUpdate, error)`
Returns comprehensive status information including position, velocity, current, and status flags.

#### `DCStatusUpdateToSI(update DCStatusUpdate) DCStatusUpdateSI`
Converts raw status data to SI units (millimeters, mm/s) based on the configured stage and motor types.

```go
status, err := controller.GetDCStatusUpdate(1)
statusSI := controller.DCStatusUpdateToSI(status)

fmt.Printf("Position: %.3f mm\n", statusSI.Position)
fmt.Printf("Velocity: %.3f mm/s\n", statusSI.Velocity)
fmt.Printf("Is Homed: %t\n", statusSI.StatusBits.IsHomed)
fmt.Printf("In Motion: %t\n", statusSI.StatusBits.InMotionCW || statusSI.StatusBits.InMotionCCW)
```

## Data Types

### Direction
- `Forward` - Move in forward direction
- `Reverse` - Move in reverse direction

### StopMode
- `Abrupt` - Immediate stop
- `Soft` - Gradual deceleration stop

### VelocityProfile
Structure containing minimum velocity, maximum velocity, and acceleration parameters.

### JogParameters
Structure containing jog mode, step size, velocities, acceleration, and stop mode.

### DCStatusBits
Comprehensive status flags including:
- Hard and soft limit states
- Motion direction indicators
- Homing and initialization status
- Error conditions and faults
- Power and enable states

## Unit Conversions

The library automatically handles conversions between physical units (millimeters, mm/s, mm/s²) and internal controller counts based on the specified stage and motor types. All public APIs use real-world units for ease of use.

### Internal Conversion Methods
- `PositionToCounts(position float64) int32`
- `CountsToPosition(counts int32) float64`
- `VelocityToCounts(velocity float64) uint32`
- `CountsToVelocity(counts uint32) float64`
- `AccelerationToCounts(acceleration float64) uint32`
- `CountsToAcceleration(counts int32) float64`

## Error Handling

The library provides specific error constants:
- `ErrChannelNotSupported` - Invalid channel number (KDC101 only supports channel 1)

Standard Go error handling patterns apply for communication errors, invalid parameters, and hardware faults.

## Thread Safety

This library is **not** thread-safe. If you need to control the same controller from multiple goroutines, you must implement your own synchronization mechanisms.

## License

This project is authored by Leonardo Rossi Leao and was created on September 26th, 2025.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
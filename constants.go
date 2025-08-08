package sphero

const (
	DataPacketStart = 0x8D
	DataPacketEnd   = 0xD8

	FlagIsResponse                = 0x01
	FlagRequestsResponse          = 0x02
	FlagRequestsOnlyErrorResponse = 0x04
	FlagResetsInactivityTimeout   = 0x08

	DevicePowerInfo    = 0x13
	DeviceUserIO       = 0x1a
	DeviceDriving      = 0x16
	DeviceAPIProcessor = 0x10
	DeviceSystemInfo   = 0x11
	DeviceAnimatronics = 0x17
	DeviceSensor       = 0x18

	PowerCommandsDeepSleep      = 0x00
	PowerCommandsSleep          = 0x01
	PowerCommandsBatteryVoltage = 0x03
	PowerCommandsWake           = 0x0D

	SensorCommandMask                   = 0x00
	SensorCommandResponse               = 0x02
	SensorCommandConfigureCollision     = 0x11
	SensorCommandCollisionDetectedAsync = 0x12
	SensorCommandResetLocator           = 0x13
	SensorCommandEnableCollisionAsync   = 0x14
	SensorCommandSensor1                = 0x0f
	SensorCommandSensor2                = 0x17
	SensorCommandConfigureSensorStream  = 0x0c

	UserIOCommandsAllLEDs = 0x0e

	DrivingCommandsWithHeading = 0x07

	SystemInfoCommandsBootLoaderVersion = 0x01
)

// CollisionConfig provides configuration for the collision detection alogorithm.
// For more information refer to the api specification of "Orbotix Communication API"
// see also: http://wiki.mark-toma.com/view/Sphero_API_Tutorial
type CollisionConfig struct {
	// Detection method type to use. Methods 01h and 02h are supported as
	// of FW ver 1.42. Use 00h to completely disable this service.
	Method uint8
	// An 8-bit settable threshold for the X (left/right) axes of Sphero.
	// A value of 00h disables the contribution of that axis.
	Xt uint8
	// An 8-bit settable speed value for the X axes. This setting is ranged
	// by the speed, then added to Xt to generate the final threshold value.
	Xs uint8
	// An 8-bit settable threshold for the Y (front/back) axes of Sphero.
	// A value of 00h disables the contribution of that axis.
	Yt uint8
	// An 8-bit settable speed value for the Y axes. This setting is ranged
	// by the speed, then added to Yt to generate the final threshold value.
	Ys uint8
	// An 8-bit post-collision dead time to prevent retriggering; specified
	// in 10ms increments.
	Dead uint8
}

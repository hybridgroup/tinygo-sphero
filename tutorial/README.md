# Tutorial

This tutorial contains a series of small activities to help you learn how to control the Sphero Mini robot using Bluetooth.

## Finding the MAC address or Bluetooth ID for the Sphero robot

You will need to determine what the MAC address or Bluetooth ID is for the Sphero Mini robot you want to connect to.

On Linux and Windows you will use the MAC address of the device to connect.

On macOS you must use the Bluetooth ID of the device to connect.

Therefore, you must know the correct name and then MAC address or ID for that device in order to connect to it.

First, install the Bluetooth scanner command:

```shell
go install ./cmd/blescanner
```

Then, run the command:

```shell
blescanner
```

It should show the names of the various Bluetooth devices around you, including the Sphero Mini you want to connect to.

## The tutorial

The tutorial steps can be run either on your computer, or on a Bluetooth enabled microcontroller such as the Pimoroni Badger2040-W.

### step1

This first step tests that the Sphero Mini is connected correctly to your computer by turning on the LED.

#### Running on your computer

```shell
go run ./tutorial/step1/ [MAC address or Bluetooth ID]
```

#### Running on your microcontroller

```shell
tinygo flash -target badger2040-w -ldflags="-X main.DeviceAddress=[MAC address]" ./tutorial/step1/
```

### step2

Now let's make the Sphero Mini roll forwards and then backwards.

#### Running on your computer

```shell
go run ./tutorial/step2/ [MAC address or Bluetooth ID]
```

#### Running on your microcontroller

```shell
tinygo flash -target badger2040-w -ldflags="-X main.DeviceAddress=[MAC address]" ./tutorial/step2/
```

### step3

Let's make the Sphero Mini roll in a square pattern then change the LED to blue when finished.

#### Running on your computer

```shell
go run ./tutorial/step3/ [MAC address or Bluetooth ID]
```

#### Running on your microcontroller

```shell
tinygo flash -target badger2040-w -ldflags="-X main.DeviceAddress=[MAC address]" ./tutorial/step3/
```

### step4

Let's make the Sphero Mini roll in a triangle pattern and change the LED color for each side.

#### Running on your computer

```shell
go run ./tutorial/step4/ [MAC address or Bluetooth ID]
```

#### Running on your microcontroller

```shell
tinygo flash -target badger2040-w -ldflags="-X main.DeviceAddress=[MAC address]" ./tutorial/step4/
```

### step5

We can also receive data notifications from the Sphero Mini such as the battery voltage. This program runs for 15 seconds and displays the current battery charge.

#### Running on your computer

```shell
go run ./tutorial/step5/ [MAC address or Bluetooth ID]
```

#### Running on your microcontroller

```shell
tinygo flash -target badger2040-w -ldflags="-X main.DeviceAddress=[MAC address]" ./tutorial/step5/
```

### step6

There are some interesting data notifications from the Sphero Mini such as collision detection. Pick up the Sphero Mini and shake it around!

#### Running on your computer

```shell
go run ./tutorial/step6/ [MAC address or Bluetooth ID]
```

#### Running on your microcontroller

```shell
tinygo flash -target badger2040-w -ldflags="-X main.DeviceAddress=[MAC address]" ./tutorial/step6/
```

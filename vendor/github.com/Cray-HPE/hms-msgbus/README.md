# HMS Message Bus Package

## Overview

This Go package encapsulates the use of the Kafka message bus system,
providing a simpler, more convenient interface for the user.  Encapsulation
of the underlying Kafka interface also allows for easy changing of the
underlying more complex Go packages if desired without changing the user's
interface.

## Usage

In general, the flow of events is:

* Create a message bus handle, for either reading or writing of messages
* For read operations, if asynchronous/control loop operation is desired, register a callback for inbound data.
* When a message is to be written, call the appropriate *write* function.
* If synchronous reads are to be done, call the appropriate *read* function.



## Interface

```
type MsgBusIO interface {
    Disconnect() error
    MessageWrite(msg string) error
    MessageRead() (string,error)
    MessageAvailable() int          //check for availability
    RegisterCB(cbfunc CBFunc) error
    UnregisterCB() error
    Status() int
}
```

## Configuration data and constants

```
type MsgBusTech int
type BlockingMode int
type BusStatus int
type BusDir int
type SubscriptionToken int
type CallbackToken int
type CBFunc func(msg string)

const (
    BusTechKafka MsgBusTech = 1
    //Add more if need be
)

const (
    NonBlocking BlockingMode = 1
    Blocking    BlockingMode = 2
)

const (
    StatusOpen   BusStatus = 1
    StatusClosed BusStatus = 2
)

const (
    BusWriter BusDir = 1
    BusReader BusDir = 2
)

// Used to configure a msgbus connection.  

type MsgBusConfig struct {
    BusTech MsgBusTech      //currently only BusTechKafka
    Host string             //msgbus host defaults to "localhost"
    Port int                //msgbus port
    Blocking BlockingMode   //Defaults to Blocking
    Direction BusDir        //BusWriter or BusReader
    ConnectRetries int      //# of times to attempt initial connection 
    Topic string            //Topic to subscribe to or inject into
}
```

## Methods

```
/////////////////////////////////////////////////////////////////////////////
// Swap out default logger with a different instance of a Logrus logger.
// Library will use default logger if this function is never called.
//
// loggerP(in): New logger to use
// Return:      None.
/////////////////////////////////////////////////////////////////////////////

SetLogger(loggerP *logrus.Logger)


/////////////////////////////////////////////////////////////////////////////
// Connect to a message bus.  This function is not part of the interface,
// but will return the correct struct for the given interface.
//
// NOTE: this is one-way at a time for now.  Bus connection can be a reader
// or a writer, but not both.  That is a potential future enhancement.  If
// an application needs both, 2 connections must be made.
//
// cfg(in):  Connection parameters.
// Return:   MsgBusIO interface, and error status
/////////////////////////////////////////////////////////////////////////////

Connect(cfg MsgBusConfig) (MsgBusIO,error)

/////////////////////////////////////////////////////////////////////////////
// Disconnect from message bus.
//
// Args:   None
// Return: Error data if there was an error.
/////////////////////////////////////////////////////////////////////////////

Disconnect() error

/////////////////////////////////////////////////////////////////////////////
// Interface function to write a message to the kafka bus.  The topic was
// determined when the connection was opened.  
//
// msg(in):  Message string to write.
// Return:   Error status of the operation.
/////////////////////////////////////////////////////////////////////////////

MessageWrite(msg string) error

/////////////////////////////////////////////////////////////////////////////
// Interface method for blocking-read of inbound messagess.  Illegal for
// writers, but needed to fill out the interface specification.
//
// Return:  Empty string, Error -- this is illegal, just a place holder.
/////////////////////////////////////////////////////////////////////////////

MessageRead() (string,error)

/////////////////////////////////////////////////////////////////////////////
// Interface function for checking for message availability on a reader
// connection.  Note that the behavior of this function is undefined if
// a callback function has been registered; it is only valuable to use
// for blocking read connections.
//
//Return: 0 if no messages are available, > 0 == message(s) available.
/////////////////////////////////////////////////////////////////////////////

MessageAvailable() int

/////////////////////////////////////////////////////////////////////////////
// Interface method for registering a received-message callback function.
// This function is called whenever a message arrives from the subscribed
// topic.  This allows for 'event driven' programming models.
//
// NOTE: only valid for reader connections.
//
// cbfunc(in):  Function to call when messages arrive.
// Return:      Error status of operation.
/////////////////////////////////////////////////////////////////////////////

RegisterCB(cbfunc CBFunc) error

/////////////////////////////////////////////////////////////////////////////
// Interface method for un-registering a callback function.  Only valid with
// reader connections.
//
// Return: Error status of operation. (currently can't fail)
/////////////////////////////////////////////////////////////////////////////

UnregisterCB() error

/////////////////////////////////////////////////////////////////////////////
// Interface method for checking msgbus connection status.
//
// Return: Status -- open or closed.
/////////////////////////////////////////////////////////////////////////////

Status() int
```


## Use Cases And Examples

```
/////////////////////////////////////////////////////////////////////////////
//Opening a connection to a message bus for writing messages:
/////////////////////////////////////////////////////////////////////////////

...
    // Message bus connection configuration

    mcfg := msgbus.MsgBusConfig{BusTech: msgbus.BusTechKafka,
                                Host: "localhost",
                                Port: 9092,
                                Blocking: msgbus.Blocking,
                                Direction: msgbus.BusWriter,
                                ConnectRetries: 10,
                                Topic: "hb_events",}

    // Connect to message bus

    mbusW,err := msgbus.Connect(mcfg)

    if (err != nil) {
        log.Printf("Error connecting to bus: %v",err)
        return
    }

...
```

```
/////////////////////////////////////////////////////////////////////////////
//Writing a message:
/////////////////////////////////////////////////////////////////////////////

...
    msg := fmt.Sprintf("%d inches of rain in Spain fell mainly on the plain.",
				inchesTotal)

    //Will potentially block if the connection was opened in Blocking mode;
    //if opened in NonBlocking mode will not block.

    err := mbusW.MessageWrite(msg)
    if (err != nil) {
        log.Printf("ERROR writing message to bus: %v",err)
    }
...
```

```
/////////////////////////////////////////////////////////////////////////////
//Opening a connection to a message bus for reading messages:
/////////////////////////////////////////////////////////////////////////////

...
    // Message bus connection configuration

    mcfg := msgbus.MsgBusConfig{BusTech: msgbus.BusTechKafka,
                                Host: "localhost",
                                Port: 9092,
                                Blocking: msgbus.Blocking,
                                Direction: msgbus.BusReader,
                                ConnectRetries: 10,
                                Topic: "hb_events",}

    // Connect to message bus

    mbusR,err := msgbus.Connect(mcfg)

    if (err != nil) {
        log.Printf("Error connecting to bus: %v",err)
        return
    }

...
```

```
/////////////////////////////////////////////////////////////////////////////
//Reading a message using a blocking read operation:
/////////////////////////////////////////////////////////////////////////////

...
    //First check if a message is available (if desired)

    if (mbusR.MessageAvailable() != 0) {
        msg,err := mbusR.MessageRead() //will block until message is ready
        if (err != nil) {
            log.Printf("ERROR reading message: %v",err)
        } else {
            log.Printf("Message received: '%s'",msg)
        }
    }
...
```

```
/////////////////////////////////////////////////////////////////////////////
//Reading a message using a callback function:
/////////////////////////////////////////////////////////////////////////////

func my_cbfunc(msg string) {
    log.Printf("Message Received: '%s'",msg)
}

func myfunc() {
...
    // Open a message bus reader connection, as above

    // Message bus connection configuration

    mcfg := msgbus.MsgBusConfig{BusTech: msgbus.BusTechKafka,
                                Host: "localhost",
                                Port: 9092,
                                Blocking: msgbus.Blocking,
                                Direction: msgbus.BusReader,
                                ConnectRetries: 10,
                                Topic: "hb_events",}

    // Connect to message bus

    mbusR,err := msgbus.Connect(mcfg)

    if (err != nil) {
        log.Printf("Error connecting to bus: %v",err)
        return
    }

    //Register a function to be called when messages arrive

    err = mbusR.RegisterCB(my_cbfunc)
    if (err != nil) {
        log.Printf("ERROR registering callback function: %v",err)
    }

    //Do other stuff.  my_cbfunc() fill be called when messages arrive.

    ...

    //Un-register callback function if you no longer want to use it.

    mbusR.UnregisterCB()

    ...
```

```
/////////////////////////////////////////////////////////////////////////////
//Closing a connection
/////////////////////////////////////////////////////////////////////////////

...
   mbusW.Disconnect()
...
```

## NOTES

* At this time a connection can be opened for reading or writing, but not both.   If both reading and writing are to be done, separate connections and handles must be used.

* There is currently no re-connect logic.  If a connection dies, the application will have to close and then re-open a new connection.

* Trying to use a callback function and MessageAvailable() and MessageRead() will result in undefined behavior.


**BUILD NOTES:**

* Any microservice using hms-msgbus starting at version 1.11.0 needs to use '-tags musl' in any 'go build' or 'go test' instruction, which will affect Dockerfiles.



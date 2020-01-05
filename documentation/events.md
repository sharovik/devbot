##Events
This feature will help you to improve the skills of your bot. Now you are able to create your own event for your custom message.

### How to create a new event
As there is a million possible events which users can create, we introduced a file `defined-events.go` where you can initialise your custom events. 

* run `cp defined-events.go.dist defined-events.go` to create the file where you will define your events
* create your event directory in `events` directory. Ex: `events/my-brand-new-event`
* create in your new directory file with name `my-event`. There is no black magic inside the naming, we just introduce the structured way of how to define the event files.
* create the logic for your new event struct object and make sure that this logic is compatible with the interface `events.Event`, which you can find this interface here `events/1/main-event.go:3`
* add your object to the "map" of the events `events.DefinedEvents` in init method of defined-events.go file 
```DefinedEvents.Events["CUSTOM_EVENT"] = your_package.Event```
* add to the dictionary, message regex by which your event will be triggered
* have fun

### Example of the event
You can find it here ```events/themer-wordpress```
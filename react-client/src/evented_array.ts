class EventedArray {
  public Stack: any[]
  private handler: () => void

  constructor(handler: () => void) {
    console.log("Handler: ", handler)
    this.handler = handler
    this.Stack = []
  }

  public callHandler() {
    console.log("This: ", this)
    this.handler();
  }

  public push(o: any) {
      this.Stack.push(o);
      this.callHandler();
  }
}

export default EventedArray
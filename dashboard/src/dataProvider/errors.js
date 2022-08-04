export class NotImplementedError extends Error {
    constructor(message) {
      super(message);
  
      this.message = message;
      this.name = 'NotImplementedError';
    }
  }
  
  export class HttpError extends Error {
    constructor(e, status) {
      super(e.data.message);

      this.message = e.data.message;
      this.status = status;
      this.name = 'HttpError';
    }
  } 
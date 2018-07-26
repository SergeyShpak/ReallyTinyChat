import * as React from 'react';
import * as Props from './props';

class ErrorBoundary extends React.Component<Props.IErrorBoundaryProps, {
  hasError: boolean,
  onError: () => void
}> {
  constructor(props: Props.IErrorBoundaryProps) {
    super(props);
    this.state = 
    {
      hasError: false,
      onError: props.onError,
    };
  }

  public componentDidCatch(error, info) {
    this.setState({ hasError: true });
  }

  public render() {
    if (this.state.hasError) {
      this.state.onError()
    }
    return this.props.children;
  }
}

export default ErrorBoundary
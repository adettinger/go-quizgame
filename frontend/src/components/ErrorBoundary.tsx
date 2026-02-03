import { Component, type ErrorInfo, type ReactNode } from 'react';
import logger from '../utils/Logger';

interface Props {
    children: ReactNode;
    fallback?: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

class ErrorBoundary extends Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = { hasError: false, error: null };
    }

    static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
        logger.error(`React error boundary caught an error: ${error.message}`, {
            stack: error.stack,
            componentStack: errorInfo.componentStack
        });
    }

    render(): ReactNode {
        if (this.state.hasError) {
            return this.props.fallback || (
                <div className="error-container p-4 bg-light border rounded">
                    <h2 className="text-danger">Something went wrong</h2>
                    <p>Please try again later or contact support if the problem persists.</p>
                    <button
                        className="btn btn-primary mt-3"
                        onClick={() => this.setState({ hasError: false, error: null })}
                    >
                        Try again
                    </button>
                </div>
            );
        }

        return this.props.children;
    }
}

export default ErrorBoundary;
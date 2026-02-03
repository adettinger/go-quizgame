enum LogLevel {
    DEBUG = 'DEBUG',
    INFO = 'INFO',
    WARN = 'WARN',
    ERROR = 'ERROR'
}

class Logger {
    private appName: string;

    constructor(appName: string = 'Frontend') {
        this.appName = appName;

        // Set up global error handlers
        this.setupGlobalErrorHandlers();
    }

    private formatMessage(level: LogLevel, message: string): string {
        const timestamp = new Date().toISOString();
        return `[${timestamp}] [${this.appName}] [${level}] ${message}`;
    }

    private setupGlobalErrorHandlers(): void {
        // Handle uncaught exceptions
        window.addEventListener('error', (event) => {
            this.error(`Uncaught error: ${event.message}`, {
                filename: event.filename,
                lineno: event.lineno,
                colno: event.colno,
                stack: event.error?.stack
            });
        });

        // Handle unhandled promise rejections
        window.addEventListener('unhandledrejection', (event) => {
            const reason = event.reason?.toString() || 'Unknown reason';
            this.error(`Unhandled promise rejection: ${reason}`, {
                stack: event.reason?.stack
            });
        });
    }

    public debug(message: string, data?: any): void {
        const formattedMessage = this.formatMessage(LogLevel.DEBUG, message);
        if (data) {
            console.debug(formattedMessage, data);
        } else {
            console.debug(formattedMessage);
        }
    }

    public info(message: string, data?: any): void {
        const formattedMessage = this.formatMessage(LogLevel.INFO, message);
        if (data) {
            console.info(formattedMessage, data);
        } else {
            console.info(formattedMessage);
        }
    }

    public warn(message: string, data?: any): void {
        const formattedMessage = this.formatMessage(LogLevel.WARN, message);
        if (data) {
            console.warn(formattedMessage, data);
        } else {
            console.warn(formattedMessage);
        }
    }

    public error(message: string, data?: any): void {
        const formattedMessage = this.formatMessage(LogLevel.ERROR, message);
        if (data) {
            console.error(formattedMessage, data);
        } else {
            console.error(formattedMessage);
        }
    }
}

// Create a singleton instance
const logger = new Logger();
export default logger;
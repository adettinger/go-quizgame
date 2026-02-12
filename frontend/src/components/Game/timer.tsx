import { useState, useEffect } from 'react';
import { differenceInSeconds } from 'date-fns';
import { Card, Flex, Text } from '@radix-ui/themes';

interface CountdownTimerProps {
    deadline: string | Date; // ISO string or Date object
    onExpire?: () => void;
}

export function CountdownTimer(
    props: CountdownTimerProps,
) {
    const [secondsLeft, setSecondsLeft] = useState<number>(0);
    const [isExpired, setIsExpired] = useState<boolean>(false);

    useEffect(() => {
        const targetDate = props.deadline instanceof Date ? props.deadline : new Date(props.deadline);
        const now = new Date();

        const initialSecondsLeft = Math.max(0, differenceInSeconds(targetDate, now));
        setSecondsLeft(initialSecondsLeft);

        const intervalId = setInterval(() => {
            const now = new Date();
            const remaining = Math.max(0, differenceInSeconds(targetDate, now));

            setSecondsLeft(remaining);

            if (remaining <= 0 && !isExpired) {
                setIsExpired(true);
                if (props.onExpire) props.onExpire();
                clearInterval(intervalId);
            }
        }, 1000);

        // Clean up
        return () => clearInterval(intervalId);
    }, [props.deadline, props.onExpire, isExpired]);

    const formatTimeRemaining = () => {
        if (isExpired) return "Time's up!";

        const hours = Math.floor(secondsLeft / 3600);
        const minutes = Math.floor((secondsLeft % 3600) / 60);
        const seconds = secondsLeft % 60;

        return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    };

    return (
        <Flex gap="3" justify={'center'} align={"center"}>
            <Text>Time remaining: </Text>
            <Card style={{ backgroundColor: secondsLeft <= 10 ? "red" : "" }}>
                <Text style={{ color: "black" }}>{formatTimeRemaining()}</Text>
            </Card>
        </Flex>
    );
};
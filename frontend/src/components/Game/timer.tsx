import { useState, useEffect, useRef } from 'react';
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
    const intervalIdRef = useRef<number | null>(null);

    useEffect(() => {
        setIsExpired(false);
        const targetDate = props.deadline instanceof Date ? props.deadline : new Date(props.deadline);
        const now = new Date();

        const initialSecondsLeft = Math.max(0, differenceInSeconds(targetDate, now));
        setSecondsLeft(initialSecondsLeft);

        // If already expired on mount, do nothing
        if (initialSecondsLeft <= 0) {
            setIsExpired(true);
            // if (props.onExpire) props.onExpire();
            return; // Don't start the timer if alrea   dy expired
        }

        intervalIdRef.current = setInterval(() => {
            const now = new Date();
            const remaining = Math.max(0, differenceInSeconds(targetDate, now));

            setSecondsLeft(remaining);

            if (remaining <= 0 && !isExpired) {
                setIsExpired(true);
                if (props.onExpire) props.onExpire();

                // Clear the interval when expired
                if (intervalIdRef.current !== null) {
                    clearInterval(intervalIdRef.current);
                    intervalIdRef.current = null;
                }
            }
        }, 1000);

        // Clean up function that runs when component unmounts
        return () => {
            if (intervalIdRef.current !== null) {
                clearInterval(intervalIdRef.current);
                intervalIdRef.current = null;
            }
        };
    }, [props.deadline]);

    // useEffect(() => {

    // }, [props.onExpire]);

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
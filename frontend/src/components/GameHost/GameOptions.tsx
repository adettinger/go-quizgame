import { Button, Flex, Switch, TextField, Tooltip } from "@radix-ui/themes";
import * as Form from "@radix-ui/react-form";
import { useEffect, useState } from "react";
import { ProblemPicker } from "./ProblemPicker";
import { useProblems } from "../../hooks/useProblems";
import type { Problem } from "../../types/problem";


export function GameOptions() {
    const { data, isLoading, isError, error } = useProblems();

    const [questionTimeLimit, setQuestionTimeLimit] = useState("0");
    const [isHostAPlayer, setIsHostAPlayer] = useState(false);
    const [availableQuestions, setAvailableQuestions] = useState<Problem[]>([]);
    const [selectedQuestions, setSelectedQuestions] = useState<Problem[]>([]);

    const areAllFieldsValid = (): boolean => {
        if (
            parseInt(questionTimeLimit, 10) >= 0 &&
            selectedQuestions.length > 0
        ) {
            return true
        }
        return false;
    };


    useEffect(() => {
        if (!!data && data.length > 0 && availableQuestions.length === 0 && selectedQuestions.length === 0) {
            setAvailableQuestions(data);
        }
    }, [data]);

    const selectQuestion = (questionToAdd: Problem) => {
        setSelectedQuestions(prev => [...prev, questionToAdd]);
        setAvailableQuestions(prev => prev.filter(question => question.Id !== questionToAdd.Id));
    };

    const deselectQuestion = (questionToRemove: Problem) => {
        setSelectedQuestions(prev => prev.filter(question => question.Id !== questionToRemove.Id));
        setAvailableQuestions(prev => [...prev, questionToRemove]);
    };

    return (
        // timelimit on questions
        // is host a player
        // question list
        // ...
        <Form.Root>
            <Flex direction="column" justify={"center"} align={"center"} gap="3">
                {/* Time limit on questions */}
                <Form.Field name="timeLimit">
                    <Flex direction={"row"} gap="1">
                        <Tooltip content={"0 for no time limit"}>
                            <Form.Label>Time limit per question (seconds):</Form.Label>
                        </Tooltip>
                        <Form.Control asChild>
                            <TextField.Root
                                value={questionTimeLimit}
                                onChange={(event) => { setQuestionTimeLimit(event.target.value.replace(/[^0-9\s]/g, '')) }}
                            />
                        </Form.Control>
                    </Flex>
                </Form.Field>
                {/* is host a player */}
                <Form.Field name="isHostAPlayer">
                    <Flex direction="row" gap="3">
                        <Form.Label>Host is a player:</Form.Label>
                        <Form.Control asChild>
                            <Switch checked={isHostAPlayer} onCheckedChange={setIsHostAPlayer} />
                        </Form.Control>
                    </Flex>
                </Form.Field>

                <Form.Field name="selectedQuestions">
                    <Form.Control asChild>
                        <ProblemPicker
                            selectedQuestions={selectedQuestions}
                            availableQuestions={availableQuestions}
                            selectQuestion={selectQuestion}
                            deselectQuestion={deselectQuestion}
                            isLoading={isLoading}
                            isError={isError}
                            error={error}
                        />
                    </Form.Control>
                </Form.Field>

                {/* question list */}
                <Form.Submit asChild>
                    <Button disabled={!areAllFieldsValid()}>Start a game!</Button>
                </Form.Submit>
            </Flex>
        </Form.Root>
    );
};
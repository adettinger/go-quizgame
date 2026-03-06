import { useState } from "react";
import { Button, Flex, Text, TextField, Tooltip } from "@radix-ui/themes";
import * as Form from "@radix-ui/react-form";
import { QuestionMarkCircledIcon } from "@radix-ui/react-icons";
import { ProblemPicker } from "./ProblemPicker";

interface GameHostFormProps {
    onSubmit: (timeLimit: number, questionIds: string[]) => void;
};


export function GameHostForm({ onSubmit }: GameHostFormProps) {

    const [questionTimeLimit, setQuestionTimeLimit] = useState("0");
    const [selectedQuestionIds, setSelectedQuestionIds] = useState<string[]>([]);

    const areAllFieldsValid = (): boolean => {
        if (
            parseInt(questionTimeLimit, 10) > 0 &&
            selectedQuestionIds.length > 0
        ) {
            return true
        }
        return false;
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onSubmit(parseInt(questionTimeLimit, 10), selectedQuestionIds);
    };

    return (
        <Form.Root onSubmit={handleSubmit} onKeyDown={(e) => e.key === 'Enter' && e.preventDefault()}>
            <Flex direction={"column"} justify={"center"} align={"center"} gap="3">

                {/* Time limit on questions */}
                <Form.Field name="timeLimit">
                    <Flex direction={"row"} gap="1">
                        <Form.Label>Time limit per question (seconds):</Form.Label>
                        <Tooltip content={"0 for no time limit"}>
                            <QuestionMarkCircledIcon />
                        </Tooltip>
                        <Form.Control asChild>
                            <TextField.Root
                                value={questionTimeLimit}
                                onChange={(event) => { setQuestionTimeLimit(event.target.value.replace(/[^0-9\s]/g, '')) }}
                            />
                        </Form.Control>
                    </Flex>
                </Form.Field>
                {/* question list */}
                <Form.Field name="selectedQuestions">
                    <Form.Control asChild>
                        <ProblemPicker
                            setQuestionIds={setSelectedQuestionIds}
                        />
                    </Form.Control>
                </Form.Field>

                <Form.Submit asChild>
                    <Button disabled={!areAllFieldsValid()}>Start a game!</Button>
                </Form.Submit>
            </Flex>
        </Form.Root>
    );
};
import { Button } from "@radix-ui/themes";
import { Form } from "radix-ui";


export function CreateProblemForm() {
    return (
        <Form.Root>
            <Form.Field name="Question">
                <Form.Label>Question</Form.Label>
                <Form.Message match="valueMissing">
                    Please enter a question
                </Form.Message>
                <Form.Control asChild>
                    <input type="text" required />
                </Form.Control>
            </Form.Field>

            <Form.Field name="Answer">
                <Form.Label>Question</Form.Label>
                <Form.Message match="valueMissing">
                    Please enter a answer
                </Form.Message>
                <Form.Control asChild>
                    <input type="text" required />
                </Form.Control>
            </Form.Field>
            <Form.Submit asChild>
                <Button>Post a question</Button>
            </Form.Submit>
        </Form.Root>
    );
};
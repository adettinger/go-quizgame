import { CaretDownIcon, CaretUpIcon } from "@radix-ui/react-icons"
import { Button, Flex, Text } from "@radix-ui/themes";
import { useState, type ReactNode } from "react";

export interface DropdownProps {
    title: string;
    children: ReactNode;
}

export function Dropdown({ title, children }: DropdownProps) {
    const [isVisible, setIsVisible] = useState(false);

    return (
        <>
            <Flex align="center" justify="between" width="100%">
                <div style={{ width: '24px' }} /> {/* Spacer to balance the layout */}
                <Text align="center" weight="bold">{title}</Text>
                <Button onClick={() => setIsVisible(!isVisible)} variant="outline" color="gray">
                    {isVisible ?
                        <CaretUpIcon fontWeight={"bold"} />
                        :
                        <CaretDownIcon fontWeight={"bold"} />
                    }
                </Button>
            </Flex>
            {isVisible && children}
        </>
    );
};
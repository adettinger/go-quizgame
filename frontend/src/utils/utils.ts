
export const ArrayToStringList = (input: string[]): string => {
    if (input.length == 0) {
        return ""
    }
    if (input.length == 1) {
        return input[0]
    }

    const lastItem = input[input.length - 1]
    const otherItems = input.slice(0, -1)

    return `${otherItems.join(", ")}, ${lastItem}`
};
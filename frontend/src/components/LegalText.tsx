import { Text } from "@radix-ui/themes";

export function LegalText() {
    return (
        <Text
            size="1"
            color="gray"
            className="legal-footer"
            style={{
                fontSize: '0.75rem',
                opacity: 0.4,
                padding: '1rem',
                textAlign: 'center',
                position: 'relative',
                zIndex: 10,
                background: 'transparent'
            }}
        >
            © 2026 Your Company. All rights to your soul are relased to me. Terms and conditions apply.
        </Text>
    );
}
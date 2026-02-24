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
                position: 'absolute',
                bottom: 0,
                left: '50%',
                transform: 'translateX(-50%)',
                zIndex: 10,
                marginTop: 'auto',
                background: 'transparent'
            }}
        >
            © 2026 Your Company. All rights to your soul are relased to me. Terms and conditions apply.
        </Text>
    );
}
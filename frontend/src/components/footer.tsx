import { Flex } from "@radix-ui/themes";
import { LegalText } from "./LegalText";
import { SocialMediaPins } from "./SocialMediaPins/SocialMediaPins";


export function Footer() {
    return (
        <Flex
            direction="row"
            justify="between"
            align="center"
            className='footer'
            style={{
                position: 'relative',
                top: '1rem',
                width: '100%',
                padding: '1rem',
                marginTop: 'auto',
            }}>
            {/* Empty div to balance the layout for centering */}
            <div style={{ flex: 1 }} />
            <div style={{ flex: 2, display: 'flex', justifyContent: 'center' }}>
                <LegalText />
            </div>
            <div style={{ flex: 1, display: 'flex', justifyContent: 'flex-end' }}>
                <SocialMediaPins />
            </div>
        </Flex>
    );
}
import React from 'react';
import { GitHubLogoIcon, LinkedInLogoIcon, TwitterLogoIcon } from '@radix-ui/react-icons';
import { Flex, IconButton, Tooltip } from '@radix-ui/themes';
import './SocialMediaPins.scss'

interface SocialLink {
    name: string;
    url: string;
    icon: React.ReactNode;
}

export function SocialMediaPins() {
    const socialLinks: SocialLink[] = [
        {
            name: 'GitHub',
            url: 'https://github.com/adettinger',
            icon: <GitHubLogoIcon width={20} height={20} />,
        },
        {
            name: 'LinkedIn',
            url: 'https://www.linkedin.com/in/alex-dettinger/',
            icon: <LinkedInLogoIcon width={20} height={20} />,
        },
        {
            name: 'Twitter',
            url: 'https://twitter.com/yourusername',
            icon: <TwitterLogoIcon width={20} height={20} />,
        },
        // Add more social media platforms as needed
    ];

    return (
        <div
            className='media-icon-container'
        >
            <Flex gap="2" direction="row">
                {socialLinks.map((social) => (
                    <Tooltip key={social.name} content={social.name}>
                        <IconButton
                            variant="soft"
                            color="gray"
                            radius="full"
                            asChild
                            size="2"
                            className='media-icon'
                        >
                            <a
                                href={social.url}
                                target="_blank"
                                rel="noopener noreferrer"
                                aria-label={`Visit my ${social.name}`}
                            >
                                {social.icon}
                            </a>
                        </IconButton>
                    </Tooltip>
                ))}
            </Flex>
        </div>
    );
};
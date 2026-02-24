// Banner.tsx
import React from 'react';
import bannerImage from '../../assets/NoyQBanner.png';
import './Banner.scss'; // We'll create this next

interface BannerProps {
    altText?: string;
}

const Banner: React.FC<BannerProps> = ({ altText = 'Website Banner' }) => {
    return (
        <div className="banner-container">
            <img
                src={bannerImage}
                alt={altText}
                className="banner-image"
                loading="lazy"
            />
        </div>
    );
};

export default Banner;
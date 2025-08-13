import React from 'react';
import { FigmaFile, HighlightableItem } from '../types/figma';

interface ImageViewerProps {
  figmaFile: FigmaFile;
  items: HighlightableItem[];
  highlightedItemId: string | null;
  onItemClick: (itemId: string) => void;
}

const ImageViewer: React.FC<ImageViewerProps> = ({
  figmaFile,
  items,
  highlightedItemId,
  onItemClick,
}) => {
  const containerStyle: React.CSSProperties = {
    position: 'relative',
    width: '100%',
    height: '100%',
    overflow: 'auto',
    padding: '20px',
  };

  const imageStyle: React.CSSProperties = {
    maxWidth: '100%',
    height: 'auto',
    display: 'block',
    border: '1px solid #ccc',
  };

  return (
    <div style={containerStyle}>
      <div style={{ position: 'relative', display: 'inline-block' }}>
        <img
          src={figmaFile.image_url}
          alt={figmaFile.name}
          style={imageStyle}
          onLoad={(e) => {
            // Store actual image dimensions for marker positioning
            const img = e.target as HTMLImageElement;
            img.dataset.actualWidth = figmaFile.canvas_width.toString();
            img.dataset.actualHeight = figmaFile.canvas_height.toString();
          }}
        />
        
        {/* Overlay markers for each item */}
        {items.map((item) => (
          <Marker
            key={item.id}
            item={item}
            isHighlighted={item.id === highlightedItemId}
            onClick={() => onItemClick(item.id)}
            imageElement={document.querySelector('img') as HTMLImageElement}
          />
        ))}
      </div>
    </div>
  );
};

interface MarkerProps {
  item: HighlightableItem;
  isHighlighted: boolean;
  onClick: () => void;
  imageElement: HTMLImageElement | null;
}

const Marker: React.FC<MarkerProps> = ({ item, isHighlighted, onClick, imageElement }) => {
  // Calculate scaled position based on actual image display size
  const getScaledPosition = () => {
    if (!imageElement) return { x: 0, y: 0, width: 0, height: 0 };
    
    const actualWidth = parseInt(imageElement.dataset.actualWidth || '0');
    const actualHeight = parseInt(imageElement.dataset.actualHeight || '0');
    const displayWidth = imageElement.offsetWidth;
    const displayHeight = imageElement.offsetHeight;
    
    const scaleX = displayWidth / actualWidth;
    const scaleY = displayHeight / actualHeight;
    
    return {
      x: item.position.x * scaleX,
      y: item.position.y * scaleY,
      width: item.position.width * scaleX,
      height: item.position.height * scaleY,
    };
  };

  const scaledPos = getScaledPosition();

  const markerStyle: React.CSSProperties = {
    position: 'absolute',
    left: `${scaledPos.x}px`,
    top: `${scaledPos.y}px`,
    width: `${scaledPos.width}px`,
    height: `${scaledPos.height}px`,
    border: `2px solid ${isHighlighted ? '#ff6b6b' : 'rgba(68, 68, 255, 0.6)'}`,
    backgroundColor: isHighlighted ? 'rgba(255, 107, 107, 0.15)' : 'rgba(68, 68, 255, 0.05)',
    cursor: 'pointer',
    boxSizing: 'border-box',
    transition: 'all 0.2s ease',
    borderRadius: '3px',
  };

  const labelStyle: React.CSSProperties = {
    position: 'absolute',
    top: '-22px',
    left: '0',
    backgroundColor: isHighlighted ? '#ff6b6b' : '#5a67d8',
    color: 'white',
    padding: '2px 8px',
    fontSize: '11px',
    borderRadius: '12px',
    whiteSpace: 'nowrap',
    maxWidth: '150px',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    fontWeight: '500',
    boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
    opacity: isHighlighted ? 1 : 0.7,
  };

  return (
    <div style={markerStyle} onClick={onClick} title={`${item.type}: ${item.name}`}>
      {isHighlighted && (
        <div style={labelStyle}>
          {item.name}
        </div>
      )}
    </div>
  );
};

export default ImageViewer;

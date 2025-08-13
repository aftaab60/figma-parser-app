import React from 'react';
import { HighlightableItem } from '../types/figma';

interface ItemListProps {
  items: HighlightableItem[];
  highlightedItemId: string | null;
  onItemClick: (itemId: string) => void;
}

const ItemList: React.FC<ItemListProps> = ({ items, highlightedItemId, onItemClick }) => {
  const components = items.filter(item => item.type === 'component');
  const instances = items.filter(item => item.type === 'instance');

  return (
    <div style={{ marginTop: '20px' }}>
      {components.length > 0 && (
        <div style={{ marginBottom: '20px' }}>
          <h4 style={{ margin: '0 0 10px 0', color: '#4444ff' }}>
            Components ({components.length})
          </h4>
          <div style={{ maxHeight: '200px', overflowY: 'auto' }}>
            {components.map(item => (
              <ItemListItem
                key={item.id}
                item={item}
                isHighlighted={item.id === highlightedItemId}
                onClick={() => onItemClick(item.id)}
              />
            ))}
          </div>
        </div>
      )}

      {instances.length > 0 && (
        <div>
          <h4 style={{ margin: '0 0 10px 0', color: '#44aa44' }}>
            Instances ({instances.length})
          </h4>
          <div style={{ maxHeight: '200px', overflowY: 'auto' }}>
            {instances.map(item => (
              <ItemListItem
                key={item.id}
                item={item}
                isHighlighted={item.id === highlightedItemId}
                onClick={() => onItemClick(item.id)}
              />
            ))}
          </div>
        </div>
      )}

      {items.length === 0 && (
        <div style={{ color: '#666', fontStyle: 'italic' }}>
          No components or instances found
        </div>
      )}
    </div>
  );
};

interface ItemListItemProps {
  item: HighlightableItem;
  isHighlighted: boolean;
  onClick: () => void;
}

const ItemListItem: React.FC<ItemListItemProps> = ({ item, isHighlighted, onClick }) => {
  const itemStyle: React.CSSProperties = {
    padding: '8px 12px',
    margin: '2px 0',
    border: '1px solid #ddd',
    borderRadius: '4px',
    backgroundColor: isHighlighted ? '#e6f3ff' : '#f9f9f9',
    cursor: 'pointer',
    transition: 'all 0.2s ease',
    borderLeft: `4px solid ${item.type === 'component' ? '#4444ff' : '#44aa44'}`,
  };

  const nameStyle: React.CSSProperties = {
    fontWeight: 'bold',
    marginBottom: '4px',
    color: isHighlighted ? '#0066cc' : '#333',
  };

  const positionStyle: React.CSSProperties = {
    fontSize: '12px',
    color: '#666',
  };

  return (
    <div 
      style={itemStyle} 
      onClick={onClick}
      onMouseEnter={(e) => {
        e.currentTarget.style.backgroundColor = isHighlighted ? '#d4edda' : '#f0f0f0';
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.backgroundColor = isHighlighted ? '#e6f3ff' : '#f9f9f9';
      }}
    >
      <div style={nameStyle}>{item.name}</div>
      <div style={positionStyle}>
        Position: ({Math.round(item.position.x)}, {Math.round(item.position.y)})
        <br />
        Size: {Math.round(item.position.width)} × {Math.round(item.position.height)}
      </div>
    </div>
  );
};

export default ItemList;

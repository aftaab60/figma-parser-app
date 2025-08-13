import React, { useState, useEffect } from 'react';
import { FigmaFile, FigmaComponent, FigmaInstance, HighlightableItem, FigmaFileDetails } from '../types/figma';
import ImageViewer from '../components/ImageViewer';
import ItemList from '../components/ItemList';
import { api } from '../utils/api';

const HomePage: React.FC = () => {
  const [figmaFile, setFigmaFile] = useState<FigmaFile | null>(null);
  const [components, setComponents] = useState<FigmaComponent[]>([]);
  const [instances, setInstances] = useState<FigmaInstance[]>([]);
  const [highlightedItemId, setHighlightedItemId] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const items: HighlightableItem[] = [
    ...components.map(comp => ({
      id: comp.id.toString(),
      name: comp.name,
      type: 'component' as const,
      position: { x: comp.x, y: comp.y, width: comp.width, height: comp.height }
    })),
    ...instances.map(inst => ({
      id: inst.id.toString(),
      name: inst.name,
      type: 'instance' as const,
      position: { x: inst.x, y: inst.y, width: inst.width, height: inst.height }
    }))
  ];

  const handleItemClick = (itemId: string) => {
    setHighlightedItemId(itemId === highlightedItemId ? null : itemId);
  };

  const handleFigmaUrlSubmit = async (figmaUrl: string, figmaToken: string) => {
    setLoading(true);
    setError(null);
    
    try {
      // Step 1: Parse the Figma file with the token
      const parseResult = await api.parseFigmaFile(figmaUrl, figmaToken);
      console.log('Parse result:', parseResult);
      
      // Step 2: Get the complete file details with components and instances
      const fileDetails = await api.getFigmaFileDetails(parseResult.data.id);
      console.log('File details:', fileDetails);
      
      // Step 3: Set the state with the fetched data
      const detailsData: FigmaFileDetails = fileDetails.data;
      setFigmaFile(detailsData.file);
      setComponents(detailsData.components || []);
      setInstances(detailsData.instances || []);
      
    } catch (err) {
      console.error('Error:', err);
      if (err instanceof Error) {
        // Handle specific error messages from the backend
        if (err.message.includes('Invalid Figma token')) {
          setError('Invalid Figma token. Please check your personal access token.');
        } else if (err.message.includes('insufficient permissions')) {
          setError('Insufficient permissions. Make sure you have access to this Figma file.');
        } else {
          setError(err.message);
        }
      } else {
        setError('An unexpected error occurred');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', height: '100vh' }}>
      {/* Left Panel - Controls and Item List */}
      <div style={{ width: '300px', padding: '20px', borderRight: '1px solid #ccc' }}>
        <h1>Figma Parser Admin</h1>
        
        {/* Figma URL Input */}
        <FigmaUrlInput onSubmit={handleFigmaUrlSubmit} loading={loading} />
        
        {error && (
          <div style={{ color: 'red', margin: '10px 0' }}>
            Error: {error}
          </div>
        )}
        
        {/* Items List */}
        {figmaFile && (
          <>
            <h3>File: {figmaFile.name}</h3>
            <ItemList 
              items={items}
              highlightedItemId={highlightedItemId}
              onItemClick={handleItemClick}
            />
          </>
        )}
      </div>
      
      {/* Right Panel - Image Viewer */}
      <div style={{ flex: 1 }}>
        {figmaFile ? (
          <ImageViewer
            figmaFile={figmaFile}
            items={items}
            highlightedItemId={highlightedItemId}
            onItemClick={handleItemClick}
          />
        ) : (
          <div style={{ 
            display: 'flex', 
            alignItems: 'center', 
            justifyContent: 'center', 
            height: '100%',
            color: '#666'
          }}>
            Upload a Figma file to get started
          </div>
        )}
      </div>
    </div>
  );
};

// Simple input component for Figma URL and Token
const FigmaUrlInput: React.FC<{ 
  onSubmit: (url: string, token: string) => void; 
  loading: boolean; 
}> = ({ onSubmit, loading }) => {
  const [url, setUrl] = useState('');
  const [token, setToken] = useState(() => {
    // Load saved token from localStorage
    return localStorage.getItem('figma_token') || '';
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (url.trim() && token.trim()) {
      // Save token to localStorage for future use
      localStorage.setItem('figma_token', token.trim());
      onSubmit(url.trim(), token.trim());
    }
  };

  const clearToken = () => {
    setToken('');
    localStorage.removeItem('figma_token');
  };

  return (
    <form onSubmit={handleSubmit} style={{ marginBottom: '20px' }}>
      <div style={{ marginBottom: '10px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <label>Figma API Token:</label>
          {token && (
            <button
              type="button"
              onClick={clearToken}
              style={{
                background: 'none',
                border: 'none',
                color: '#666',
                cursor: 'pointer',
                fontSize: '12px',
                textDecoration: 'underline'
              }}
            >
              Clear
            </button>
          )}
        </div>
        <input
          type="password"
          value={token}
          onChange={(e) => setToken(e.target.value)}
          placeholder="figd_..."
          style={{ 
            width: '100%', 
            padding: '8px', 
            marginTop: '5px',
            border: '1px solid #ccc',
            borderRadius: '4px'
          }}
          disabled={loading}
        />
        <small style={{ color: '#666', fontSize: '12px' }}>
          Get your token from Figma Settings → Personal Access Tokens
        </small>
      </div>
      
      <div style={{ marginBottom: '10px' }}>
        <label>Figma File URL:</label>
        <input
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://www.figma.com/file/..."
          style={{ 
            width: '100%', 
            padding: '8px', 
            marginTop: '5px',
            border: '1px solid #ccc',
            borderRadius: '4px'
          }}
          disabled={loading}
        />
      </div>
      
      <button 
        type="submit" 
        disabled={loading || !url.trim() || !token.trim()}
        style={{
          padding: '8px 16px',
          backgroundColor: '#007bff',
          color: 'white',
          border: 'none',
          borderRadius: '4px',
          cursor: loading ? 'not-allowed' : 'pointer',
          opacity: (loading || !url.trim() || !token.trim()) ? 0.6 : 1
        }}
      >
        {loading ? 'Parsing...' : 'Parse Figma File'}
      </button>
    </form>
  );
};

export default HomePage;

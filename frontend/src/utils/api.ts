const API_BASE_URL = 'http://localhost:3000';

export const api = {
  // Parse and save a Figma file
  parseFigmaFile: async (figmaUrl: string, figmaToken: string) => {
    const response = await fetch(`${API_BASE_URL}/parse-figma-file`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${figmaToken}`,
      },
      body: JSON.stringify({ "figma_file_url": figmaUrl }),
    });
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || `Failed to parse Figma file: ${response.statusText}`);
    }
    
    return response.json();
  },

  // Get complete file details with components and instances
  getFigmaFileDetails: async (fileId: number) => {
    const response = await fetch(`${API_BASE_URL}/figma-files/${fileId}`);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch file details: ${response.statusText}`);
    }
    
    return response.json();
  },
  
  // Get all Figma files (you might need to add this endpoint to your backend)
  getFigmaFiles: async () => {
    const response = await fetch(`${API_BASE_URL}/figma-files`);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch Figma files: ${response.statusText}`);
    }
    
    return response.json();
  },
  
  // Get components for a specific Figma file
  getComponents: async (figmaFileId: string) => {
    const response = await fetch(`${API_BASE_URL}/figma-files/${figmaFileId}/components`);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch components: ${response.statusText}`);
    }
    
    return response.json();
  },
  
  // Get instances for a specific Figma file
  getInstances: async (figmaFileId: string) => {
    const response = await fetch(`${API_BASE_URL}/figma-files/${figmaFileId}/instances`);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch instances: ${response.statusText}`);
    }
    
    return response.json();
  },
};

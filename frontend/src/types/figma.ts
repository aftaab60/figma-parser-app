export interface FigmaFile {
  id: number;
  name: string;
  url: string;
  file_key: string;
  image_url: string;
  thumbnails?: string;
  canvas_width: number;
  canvas_height: number;
  parsed_at: string;
  created_at: string;
  updated_at: string;
  active: boolean;
}

export interface FigmaComponent {
  id: number;
  figma_file_id: number;
  name: string;
  type: string;
  x: number;
  y: number;
  width: number;
  height: number;
  properties: Record<string, any>;
}

export interface FigmaInstance {
  id: number;
  figma_file_id: number;
  component_id: number;
  name: string;
  x: number;
  y: number;
  width: number;
  height: number;
  properties: Record<string, any>;
}

export interface FigmaFileDetails {
  file: FigmaFile;
  components: FigmaComponent[];
  instances: FigmaInstance[];
}

export interface MarkerPosition {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface HighlightableItem {
  id: string;
  name: string;
  type: 'component' | 'instance';
  position: MarkerPosition;
}

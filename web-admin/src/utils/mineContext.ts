const MINE_KEY = 'selected_mine_id';

export interface MineOption {
  id: number;
  name: string;
}

export const MINES: MineOption[] = [
  { id: 1, name: '东矿区' },
  { id: 2, name: '西矿区' },
];

export function getCurrentMineId(): number {
  const v = localStorage.getItem(MINE_KEY);
  return v ? parseInt(v, 10) : 1;
}

export function setCurrentMineId(id: number): void {
  localStorage.setItem(MINE_KEY, String(id));
}

export function getCurrentMineName(): string {
  const id = getCurrentMineId();
  return MINES.find((m) => m.id === id)?.name || `矿区#${id}`;
}

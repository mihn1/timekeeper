/**
 * Represents a column configuration for data tables
 */
export type Column = {
    /** Unique key identifying the column, should match a property in your data objects */
    key: string;
    
    /** Display title for the column header */
    title: string;
    
    /** Whether this column can be sorted */
    sortable: boolean;
    
    /** Optional function to format cell values for display */
    formatter?: (value: any) => string;
    
    /** Whether this column should be hidden on mobile devices */
    hideMobile?: boolean;
    
    /** Custom CSS classes to apply to cells in this column */
    cellClass?: string;
    
    /** Width of the column (e.g., '100px', '10%') */
    width?: string;
};

/**
 * Configuration for paginated data tables
 */
export type PaginationConfig = {
    pageSize: number;
    currentPage: number;
    totalItems: number;
    pageSizes: number[];
};
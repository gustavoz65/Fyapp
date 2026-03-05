
// Enums / Union Types


export type UserRole = "user" | "admin" | "premium";
export type TransactionType = "income" | "expense" | "transfer";
export type TransactionSource = "manual" | "bank_sync" | "recurring";
export type AccountType = "checking" | "savings" | "credit_card" | "investment" | "cash" | "other";
export type CategoryType = "income" | "expense";
export type BudgetPeriodType = "monthly" | "quarterly" | "yearly" | "custom";
export type GoalStatus = "in_progress" | "completed" | "cancelled";
export type NotificationType = "budget_alert" | "bill_reminder" | "goal_achieved" | "low_balance" | "transaction_alert" | "system";
export type RecurringFrequency = "daily" | "weekly" | "biweekly" | "monthly" | "quarterly" | "yearly";


// User


export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  phone?: string;
  avatar_url?: string;
  preferred_currency: string;
  preferred_language: string;
  timezone: string;
  role: UserRole;
  email_verified: boolean;
  email_verified_at?: string;
  last_login_at?: string;
  is_active: boolean;
  onboarding_completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface UserSettings {
  id: string;
  user_id: string;
  notification_email: boolean;
  notification_push: boolean;
  notification_sms: boolean;
  budget_alerts: boolean;
  bill_reminders: boolean;
  bill_reminder_days: number;
  weekly_summary: boolean;
  monthly_report: boolean;
  low_balance_alert: boolean;
  low_balance_threshold: number;
  allow_manual_transactions: boolean;
  theme: "light" | "dark" | "system";
  dashboard_layout?: string;
  created_at: string;
  updated_at: string;
}


// Bank Account


export interface BankAccount {
  id: string;
  user_id: string;
  name: string;
  bank_name?: string;
  bank_code?: string;
  account_type: AccountType;
  account_number?: string;
  agency?: string;
  initial_balance: string;
  current_balance: string;
  credit_limit?: string;
  closing_day?: number;
  due_day?: number;
  currency: string;
  color: string;
  icon: string;
  is_active: boolean;
  include_in_total: boolean;
  last_sync_at?: string;
  created_at: string;
  updated_at: string;
}


// Category


export interface Category {
  id: string;
  user_id?: string;
  name: string;
  description?: string;
  type: CategoryType;
  color: string;
  icon: string;
  is_system: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}


// Transaction


export interface Transaction {
  id: string;
  user_id: string;
  bank_account_id: string;
  category_id?: string;
  type: TransactionType;
  amount: string;
  description: string;
  notes?: string;
  source: TransactionSource;
  transaction_date: string;
  due_date?: string;
  payment_date?: string;
  is_paid: boolean;
  is_recurring: boolean;
  recurring_id?: string;
  installment_number?: number;
  total_installments?: number;
  installment_group_id?: string;
  tags?: string[];
  attachment_url?: string;
  created_at: string;
  updated_at: string;
  category?: Category;
  bank_account?: BankAccount;
}

export interface RecurringTransaction {
  id: string;
  user_id: string;
  bank_account_id: string;
  category_id?: string;
  type: TransactionType;
  amount: string;
  description: string;
  frequency: RecurringFrequency;
  day_of_month?: number;
  day_of_week?: number;
  start_date: string;
  end_date?: string;
  next_occurrence: string;
  last_generated_at?: string;
  is_active: boolean;
  auto_confirm: boolean;
  created_at: string;
  updated_at: string;
  category?: Category;
  bank_account?: BankAccount;
}


// Budget


export interface Budget {
  id: string;
  user_id: string;
  category_id?: string;
  name: string;
  amount: string;
  spent_amount: string;
  period_type: BudgetPeriodType;
  start_date: string;
  end_date: string;
  alert_threshold: string;
  alert_sent: boolean;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  category?: Category;
}

export interface BudgetSummary {
  budget: Budget;
  total_budgeted: string;
  total_spent: string;
  total_remaining: string;
  used_percentage: string;
  transaction_count: number;
}


// Goal


export interface Goal {
  id: string;
  user_id: string;
  name: string;
  description?: string;
  target_amount: string;
  current_amount: string;
  target_date?: string;
  icon: string;
  color: string;
  priority: number;
  status: GoalStatus;
  completed_at?: string;
  created_at: string;
  updated_at: string;
  contributions?: GoalContribution[];
}

export interface GoalContribution {
  id: string;
  goal_id: string;
  amount: string;
  note?: string;
  contribution_date: string;
  created_at: string;
}

export interface GoalSummary {
  total_goals: number;
  completed_goals: number;
  in_progress_goals: number;
  total_target_amount: string;
  total_saved_amount: string;
  overall_progress: string;
}


// Notification


export interface Notification {
  id: string;
  user_id: string;
  type: NotificationType;
  title: string;
  message: string;
  data?: string;
  is_read: boolean;
  read_at?: string;
  sent_via: string[];
  scheduled_for?: string;
  sent_at?: string;
  created_at: string;
}


// Dashboard


export interface DashboardSummary {
  total_balance: string;
  total_accounts: number;
  month_income: string;
  month_expense: string;
  month_balance: string;
  income_change: string;
  expense_change: string;
  active_budgets: number;
  over_budget_count: number;
  budget_usage: string;
  active_goals: number;
  goals_progress: string;
  upcoming_bills: number;
  upcoming_total: string;
  recent_transactions?: Transaction[];
  top_categories?: CategoryAmount[];
}

export interface CashFlowReport {
  period: string;
  start_date: string;
  end_date: string;
  total_income: string;
  total_expense: string;
  net_cash_flow: string;
  opening_balance: string;
  closing_balance: string;
  daily_breakdown?: DailyCashFlow[];
  income_by_category: CategoryAmount[];
  expense_by_category: CategoryAmount[];
}

export interface DailyCashFlow {
  date: string;
  income: string;
  expense: string;
  net_flow: string;
  balance: string;
}

export interface CategoryAmount {
  category_id?: string;
  category_name: string;
  amount: string;
  percentage: string;
  count: number;
}

export interface MonthlyIncomeExpense {
  month: string;
  year: number;
  income: string;
  expense: string;
  net_income: string;
  savings_rate: string;
}

export interface IncomeVsExpenseReport {
  start_date: string;
  end_date: string;
  total_income: string;
  total_expense: string;
  net_income: string;
  savings_rate: string;
  monthly_data?: MonthlyIncomeExpense[];
  income_growth: string;
  expense_growth: string;
}


// Auth


export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone?: string;
}

export interface LoginResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_at: number;
}

export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
  expires_at: number;
}

export interface SocialLoginRequest {
  provider: "google" | "facebook" | "github";
  id_token: string;
  device_info?: string;
}

export interface LinkedProvider {
  provider: string;
  email: string;
  name: string;
  avatar_url: string;
  is_primary: boolean;
  linked_at: string;
}

export interface ListProvidersResponse {
  providers: LinkedProvider[];
  has_password: boolean;
}


// Request DTOs


export interface CreateBankAccountRequest {
  name: string;
  bank_name?: string;
  bank_code?: string;
  account_type: AccountType;
  account_number?: string;
  agency?: string;
  initial_balance: string;
  credit_limit?: string;
  closing_day?: number;
  due_day?: number;
  currency?: string;
  color?: string;
  icon?: string;
  include_in_total?: boolean;
}

export interface UpdateBankAccountRequest {
  name?: string;
  bank_name?: string;
  bank_code?: string;
  account_number?: string;
  agency?: string;
  credit_limit?: string;
  closing_day?: number;
  due_day?: number;
  color?: string;
  icon?: string;
  is_active?: boolean;
  include_in_total?: boolean;
}

export interface CreateTransactionRequest {
  bank_account_id: string;
  category_id?: string;
  type: TransactionType;
  amount: string;
  description: string;
  notes?: string;
  transaction_date: string;
  due_date?: string;
  is_paid: boolean;
  tags?: string[];
  total_installments?: number;
}

export interface UpdateTransactionRequest {
  category_id?: string;
  amount?: string;
  description?: string;
  notes?: string;
  transaction_date?: string;
  due_date?: string;
  is_paid?: boolean;
  tags?: string[];
}

export interface CreateCategoryRequest {
  name: string;
  description?: string;
  type: CategoryType;
  color?: string;
  icon?: string;
}

export interface UpdateCategoryRequest {
  name?: string;
  description?: string;
  color?: string;
  icon?: string;
  is_active?: boolean;
}

export interface CreateBudgetRequest {
  category_id?: string;
  name: string;
  amount: string;
  period_type: BudgetPeriodType;
  start_date: string;
  end_date: string;
  alert_threshold?: string;
}

export interface UpdateBudgetRequest {
  name?: string;
  amount?: string;
  alert_threshold?: string;
  is_active?: boolean;
}

export interface CreateGoalRequest {
  name: string;
  description?: string;
  target_amount: string;
  target_date?: string;
  icon?: string;
  color?: string;
  priority?: number;
}

export interface UpdateGoalRequest {
  name?: string;
  description?: string;
  target_amount?: string;
  target_date?: string;
  icon?: string;
  color?: string;
  priority?: number;
  status?: GoalStatus;
}

export interface CreateGoalContributionRequest {
  amount: string;
  note?: string;
  contribution_date?: string;
}

export interface UpdateUserRequest {
  first_name?: string;
  last_name?: string;
  phone?: string;
  preferred_currency?: string;
  preferred_language?: string;
  timezone?: string;
}

export interface UpdateUserSettingsRequest {
  notification_email?: boolean;
  notification_push?: boolean;
  notification_sms?: boolean;
  budget_alerts?: boolean;
  bill_reminders?: boolean;
  bill_reminder_days?: number;
  weekly_summary?: boolean;
  monthly_report?: boolean;
  low_balance_alert?: boolean;
  low_balance_threshold?: number;
  theme?: "light" | "dark" | "system";
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

export interface SetPasswordRequest {
  new_password: string;
  confirm_password: string;
}


// Pagination


export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
  has_more: boolean;
}

export interface TransactionFilter {
  account_id?: string;
  category_id?: string;
  type?: TransactionType;
  start_date?: string;
  end_date?: string;
  is_paid?: boolean;
  min_amount?: string;
  max_amount?: string;
  search_term?: string;
  tags?: string[];
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_direction?: "asc" | "desc";
}

export interface APIError {
  code: string;
  message: string;
  details?: Record<string, string>;
}

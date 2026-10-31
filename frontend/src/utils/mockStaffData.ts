export type PipelineStage = 'Draft' | 'Quotation' | 'Confirmed' | 'Active' | 'Closed';

export interface StaffCustomer {
  id: string;
  name: string;
  contact: string;
  email: string;
  activeSubscriptions: number;
  monthlyValueCents: number;
  tier: string;
  status: 'active' | 'attention';
}

export interface QuoteLineTemplate {
  productId: string;
  product: string;
  quantity: number;
  unitPriceCents: number;
  taxRateBasisPoints: number;
}

export interface QuoteTemplate {
  id: string;
  name: string;
  plan: string;
  validityDays: number;
  lines: QuoteLineTemplate[];
}

export interface StaffQuote {
  id: string;
  customerId: string;
  customer: string;
  plan: string;
  amountCents: number;
  validUntil: string;
  status: 'Pending' | 'Approved' | 'Expired';
}

export interface StaffSubscription {
  id: string;
  customerId: string;
  customer: string;
  plan: string;
  monthlyValueCents: number;
  stage: PipelineStage;
  updatedAt: string;
}

export interface StaffInvoice {
  id: string;
  customerId: string;
  customer: string;
  invoiceNumber: string;
  issuedAt: string;
  dueDate: string;
  amountCents: number;
  status: 'Paid' | 'Outstanding';
  method?: string;
}

export interface StaffWorkspaceData {
  customers: StaffCustomer[];
  templates: QuoteTemplate[];
  quotes: StaffQuote[];
  subscriptions: StaffSubscription[];
  invoices: StaffInvoice[];
}

const customerProfiles = [
  { name: 'Northstar Analytics', contact: 'Avery Chen', email: 'avery@northstar.io' },
  { name: 'Meridian Health Group', contact: 'Jordan Patel', email: 'jordan@meridianhealth.com' },
  { name: 'Pinecone Logistics', contact: 'Morgan Ellis', email: 'morgan@pineconelogistics.co' },
  { name: 'Commonwealth Robotics', contact: 'Riley Brooks', email: 'riley@commonwealthrobotics.com' },
  { name: 'Fieldnote Commerce', contact: 'Sam Rivera', email: 'sam@fieldnote.shop' },
  { name: 'Redwood Climate', contact: 'Taylor Kim', email: 'taylor@redwoodclimate.org' },
];

const planNames = ['Starter', 'Growth', 'Scale', 'Enterprise'];
const pipelineStages: PipelineStage[] = ['Draft', 'Quotation', 'Confirmed', 'Active', 'Closed'];
const paymentMethods = ['Card', 'Bank transfer', 'Direct debit'];

const randomInt = (min: number, max: number): number =>
  Math.floor(Math.random() * (max - min + 1)) + min;

const dateOffset = (days: number): string => {
  const date = new Date();
  date.setDate(date.getDate() + days);
  return date.toISOString().slice(0, 10);
};

const makeWorkspaceData = (): StaffWorkspaceData => {
  const customers = customerProfiles.map((profile, index) => {
    const plan = planNames[randomInt(0, planNames.length - 1)];
    return {
      id: `cus_${String(index + 1).padStart(4, '0')}`,
      ...profile,
      activeSubscriptions: randomInt(1, 4),
      monthlyValueCents: randomInt(120, 2200) * 100,
      tier: `${plan} ${Math.random() > 0.5 ? 'Monthly' : 'Annual'}`,
      status: (Math.random() > 0.18 ? 'active' : 'attention') as StaffCustomer['status'],
    };
  });

  const templates: QuoteTemplate[] = [
    {
      id: 'template-growth',
      name: 'Growth onboarding',
      plan: 'Growth',
      validityDays: randomInt(21, 40),
      lines: [
        { productId: 'prod-core', product: 'Brightloop Core', quantity: 1, unitPriceCents: randomInt(330, 370) * 100, taxRateBasisPoints: 850 },
        { productId: 'prod-seat', product: 'Additional seats', quantity: 5, unitPriceCents: randomInt(22, 27) * 100, taxRateBasisPoints: 850 },
      ],
    },
    {
      id: 'template-scale',
      name: 'Scale annual',
      plan: 'Scale',
      validityDays: randomInt(21, 40),
      lines: [
        { productId: 'prod-scale', product: 'Brightloop Scale', quantity: 1, unitPriceCents: randomInt(850, 940) * 100, taxRateBasisPoints: 850 },
        { productId: 'prod-support', product: 'Priority support', quantity: 1, unitPriceCents: randomInt(110, 140) * 100, taxRateBasisPoints: 850 },
      ],
    },
    {
      id: 'template-enterprise',
      name: 'Enterprise services',
      plan: 'Enterprise',
      validityDays: randomInt(21, 40),
      lines: [
        { productId: 'prod-enterprise', product: 'Brightloop Enterprise', quantity: 1, unitPriceCents: randomInt(1800, 2050) * 100, taxRateBasisPoints: 850 },
        { productId: 'prod-implementation', product: 'Implementation package', quantity: 1, unitPriceCents: randomInt(4000, 5000) * 100, taxRateBasisPoints: 0 },
      ],
    },
  ];

  const quotes: StaffQuote[] = customers.slice(0, 4).map((customer, index) => {
    const template = templates[index % templates.length];
    const subtotal = template.lines.reduce(
      (total, line) => total + line.quantity * line.unitPriceCents,
      0
    );
    const tax = template.lines.reduce(
      (total, line) => total + Math.round(line.quantity * line.unitPriceCents * line.taxRateBasisPoints / 10000),
      0
    );
    const validityOffset = randomInt(-3, 26);
    const isExpired = validityOffset < 0 && Math.random() > 0.4;
    return {
      id: `QT-${new Date().getFullYear()}-${String(410 + index).padStart(4, '0')}`,
      customerId: customer.id,
      customer: customer.name,
      plan: template.plan,
      amountCents: subtotal + tax,
      validUntil: dateOffset(validityOffset),
      status: isExpired ? 'Expired' : Math.random() > 0.65 ? 'Approved' : 'Pending',
    };
  });

  const subscriptions: StaffSubscription[] = customers.map((customer, index) => ({
    id: `sub_${String(9800 + index).padStart(5, '0')}`,
    customerId: customer.id,
    customer: customer.name,
    plan: customer.tier,
    monthlyValueCents: customer.monthlyValueCents,
    stage: pipelineStages[(index + randomInt(0, pipelineStages.length - 1)) % pipelineStages.length],
    updatedAt: dateOffset(-randomInt(0, 18)),
  }));

  const invoices: StaffInvoice[] = customers.flatMap((customer, customerIndex) =>
    [0, 1].map((invoiceIndex) => {
      const status = (Math.random() > 0.45 ? 'Paid' : 'Outstanding') as StaffInvoice['status'];
      const daysOffset = status === 'Paid' ? -randomInt(2, 45) : randomInt(-24, 8);
      return {
        id: `inv_${customerIndex}_${invoiceIndex}`,
        customerId: customer.id,
        customer: customer.name,
        invoiceNumber: `INV-${new Date().getFullYear()}-${String(820 + customerIndex * 2 + invoiceIndex).padStart(4, '0')}`,
        issuedAt: dateOffset(daysOffset - 14),
        dueDate: dateOffset(daysOffset),
        amountCents: randomInt(17500, 265000),
        status,
        ...(status === 'Paid' ? { method: paymentMethods[randomInt(0, paymentMethods.length - 1)] } : {}),
      };
    })
  );

  return { customers, templates, quotes, subscriptions, invoices };
};

export async function fetchStaffWorkspace(): Promise<StaffWorkspaceData> {
  await new Promise<void>((resolve) => window.setTimeout(resolve, 220));
  return makeWorkspaceData();
}

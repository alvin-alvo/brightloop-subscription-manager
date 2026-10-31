import type { Invoice } from '../types';

function escapePdf(str: string): string {
  return str.replace(/\\/g, '\\\\').replace(/\(/g, '\\(').replace(/\)/g, '\\)');
}

function formatCurrency(amount: number, currency = 'USD'): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
  }).format(amount);
}

/**
 * Builds a standard, RFC-compliant PDF 1.4 document for an invoice.
 * Completely client-side with zero external binary dependencies.
 */
export function generateInvoicePdfBlob(
  invoice: Invoice,
  customerName = 'Acme Enterprises Inc.',
  customerEmail = 'billing@acme-enterprises.com'
): Blob {
  const enc = new TextEncoder();
  const dateStr = invoice.date || new Date().toISOString().split('T')[0];
  const dueDateStr = invoice.dueDate || dateStr;
  const statusStr = (invoice.status || 'open').toUpperCase();
  const currencyStr = invoice.currency || 'USD';

  const streamLines: string[] = [
    'BT',
    '/F1 20 Tf',
    '50 740 Td',
    '(BRIGHTLOOP REVENUE PLATFORM) Tj',
    '/F2 10 Tf',
    '0 -18 Td',
    '(Official B2B SaaS Tax Invoice & Settlement Receipt) Tj',
    '0 -30 Td',
    '/F1 12 Tf',
    '(INVOICE DETAILS) Tj',
    '/F2 10 Tf',
    '0 -16 Td',
    `(${escapePdf(`Invoice Number: ${invoice.invoiceNumber}`)}) Tj`,
    '0 -14 Td',
    `(${escapePdf(`Billing Date: ${dateStr}`)}) Tj`,
    '0 -14 Td',
    `(${escapePdf(`Payment Due: ${dueDateStr}`)}) Tj`,
    '0 -14 Td',
    `(${escapePdf(`Status: ${statusStr}`)}) Tj`,
    '0 -30 Td',
    '/F1 12 Tf',
    '(BILLED TO) Tj',
    '/F2 10 Tf',
    '0 -16 Td',
    `(${escapePdf(`Account: ${invoice.customerName || customerName}`)}) Tj`,
    '0 -14 Td',
    `(${escapePdf(`Email: ${invoice.customerEmail || customerEmail}`)}) Tj`,
    '0 -30 Td',
    '/F1 12 Tf',
    '(LINE ITEMS & SUMMARY) Tj',
    '/F2 10 Tf',
    '0 -16 Td',
  ];

  if (invoice.lineItems && invoice.lineItems.length > 0) {
    for (const item of invoice.lineItems) {
      streamLines.push(
        `(${escapePdf(`• ${item.description} (Qty ${item.quantity}) - ${formatCurrency(item.amount, currencyStr)}`)}) Tj`,
        '0 -14 Td'
      );
    }
  } else {
    streamLines.push(
      `(${escapePdf(`• SaaS Subscription Services - ${formatCurrency(invoice.subtotalAmount, currencyStr)}`)}) Tj`,
      '0 -14 Td'
    );
  }

  streamLines.push(
    '0 -10 Td',
    `(${escapePdf(`Subtotal: ${formatCurrency(invoice.subtotalAmount, currencyStr)}`)}) Tj`,
    '0 -14 Td',
    `(${escapePdf(`Taxes & VAT: ${formatCurrency(invoice.taxAmount, currencyStr)}`)}) Tj`,
    '0 -14 Td',
    `(${escapePdf(`Discounts Applied: -${formatCurrency(invoice.discountAmount, currencyStr)}`)}) Tj`,
    '0 -14 Td',
    '/F1 11 Tf',
    `(${escapePdf(`Total Amount: ${formatCurrency(invoice.totalAmount, currencyStr)}`)}) Tj`,
    '/F2 10 Tf',
    '0 -14 Td',
    `(${escapePdf(`Amount Paid: ${formatCurrency(invoice.amountPaid, currencyStr)}`)}) Tj`,
    '0 -14 Td',
    `(${escapePdf(`Balance Due: ${formatCurrency(invoice.amountDue, currencyStr)}`)}) Tj`,
    '0 -14 Td',
    `(${escapePdf(`Payment Method: ${invoice.paymentMethod || 'Corporate Card / Electronic Transfer'}`)}) Tj`,
    '0 -35 Td',
    '/F2 9 Tf',
    '(Thank you for choosing Brightloop. For remittance or support: billing@brightloop.io) Tj',
    'ET'
  );

  const streamText = streamLines.join('\n');
  const streamBytes = enc.encode(streamText);
  const streamLength = streamBytes.length;

  const chunks: Uint8Array[] = [];
  const offsets: number[] = [];
  let currentOffset = 0;

  function pushBytes(bytes: Uint8Array) {
    chunks.push(bytes);
    currentOffset += bytes.length;
  }

  function pushString(str: string) {
    pushBytes(enc.encode(str));
  }

  function addObj(num: number, content: string | Uint8Array) {
    offsets[num] = currentOffset;
    pushString(`${num} 0 obj\n`);
    if (typeof content === 'string') {
      pushString(content);
    } else {
      pushBytes(content);
    }
    pushString('\nendobj\n');
  }

  pushString('%PDF-1.4\n');
  addObj(1, '<< /Type /Catalog /Pages 2 0 R >>');
  addObj(2, '<< /Type /Pages /Kids [3 0 R] /Count 1 >>');
  addObj(
    3,
    '<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 4 0 R /F2 5 0 R >> >> /Contents 6 0 R >>'
  );
  addObj(4, '<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold >>');
  addObj(5, '<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>');

  // Object 6: Content stream
  const streamHeader = enc.encode(`<< /Length ${streamLength} >>\nstream\n`);
  const streamFooter = enc.encode('\nendstream');
  const fullStream = new Uint8Array(streamHeader.length + streamBytes.length + streamFooter.length);
  fullStream.set(streamHeader, 0);
  fullStream.set(streamBytes, streamHeader.length);
  fullStream.set(streamFooter, streamHeader.length + streamBytes.length);
  addObj(6, fullStream);

  const startxref = currentOffset;
  let xref = 'xref\n0 7\n0000000000 65535 f \n';
  for (let i = 1; i <= 6; i++) {
    xref += String(offsets[i]).padStart(10, '0') + ' 00000 n \n';
  }
  pushString(xref);
  pushString(`trailer\n<< /Size 7 /Root 1 0 R >>\nstartxref\n${startxref}\n%%EOF`);

  return new Blob(chunks, { type: 'application/pdf' });
}

/**
 * Triggers a native browser file download for the invoice PDF.
 */
export function downloadInvoicePdfFile(
  invoice: Invoice,
  customerName?: string,
  customerEmail?: string
): void {
  const blob = generateInvoicePdfBlob(invoice, customerName, customerEmail);
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `Invoice-${invoice.invoiceNumber || invoice.id}.pdf`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  setTimeout(() => URL.revokeObjectURL(url), 1000);
}


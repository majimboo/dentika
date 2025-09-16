# Seeded Clinics and Users

## Clinics

### 1. Dentika (ID: 1) - Main Platform
**Address:** Unit 205 JY Square Mall, Lahug, Cebu City, 6000 Cebu
**Phone:** +63 32 520 8888
**Email:** info@dentika.com
**Website:** https://dentika.com
**Tagline:** Your Trusted Dental Partner in Cebu
**Branches:** 1 (Main Branch)

#### Branch 1 - Main Branch (ID: 1)
**Address:** Unit 205 JY Square Mall, Lahug, Cebu City, 6000 Cebu
**Phone:** +63 32 520 8888
**Operating Hours:** Mon-Fri: 9:00-17:00, Sat: 9:00-12:00
**Operating Days:** Monday, Tuesday, Wednesday, Thursday, Friday, Saturday
**Open on Holidays:** No
**Closed Today:** No

**Users:**
- **Super Admin:** admin / admin (Admin User) - admin@dentika.com

**Patients:** No seeded patients (platform clinic for inventory management)

---

### 2. SmileCare Dental Clinic (ID: 2)
**Address:** G/F Ayala Center Cebu, Cebu Business Park, Cebu City, 6000 Cebu
**Phone:** +63 32 412 3456
**Email:** info@smilecare.ph
**Website:** https://smilecare.ph
**Tagline:** Caring for Your Smile with Excellence
**Branches:** 1 (Main Branch)

#### Branch 1 - Main Branch (ID: 2)
**Address:** G/F Ayala Center Cebu, Cebu Business Park, Cebu City, 6000 Cebu
**Phone:** +63 32 412 3456
**Operating Hours:** Mon-Sat: 8:00-18:00, Sun: 9:00-14:00
**Operating Days:** Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday
**Open on Holidays:** Yes
**Closed Today:** No

**Users:**
- **Admins:**
  - mcruz / admin (Maria Cruz) - maria.cruz@smilecare.ph
  - jreyes / admin (Jose Reyes) - jose.reyes@smilecare.ph
- **Doctors:**
  - asantos / admin (Ana Santos) - ana.santos@smilecare.ph
  - rgarcia / admin (Roberto Garcia) - roberto.garcia@smilecare.ph
- **Secretary:**
  - lvillanueva / admin (Luz Villanueva) - luz.villanueva@smilecare.ph
- **Assistants:**
  - cdela / admin (Carmen dela Cruz) - carmen.delacruz@smilecare.ph
  - jtanaka / admin (Juan Tanaka) - juan.tanaka@smilecare.ph

**Patients:** 20 seeded patients with Filipino names, addresses, and comprehensive medical information including:
- Personal details (name, email, phone, date of birth, gender, address)
- Emergency contacts with relationship details
- Medical information (blood type, allergies, medical conditions, current medications)
- Insurance details (PhilHealth, Maxicare, Intellicare)
- Notes for dental care considerations

---

### 3. Bright Smile Dental Center (ID: 3)
**Address:** 2/F Robinson's Galleria Cebu, Gen. Maxilom Avenue, Cebu City, 6000 Cebu
**Phone:** +63 32 234 5678
**Email:** info@brightsmile.ph
**Website:** https://brightsmile.ph
**Tagline:** Bringing Brightness to Your Smile
**Branches:** 2

#### Branch 1 - Main Branch (ID: 3)
**Address:** 2/F Robinson's Galleria Cebu, Gen. Maxilom Avenue, Cebu City, 6000 Cebu
**Phone:** +63 32 234 5678
**Operating Hours:** Mon-Fri: 9:00-19:00, Sat: 9:00-17:00
**Operating Days:** Monday, Tuesday, Wednesday, Thursday, Friday, Saturday
**Open on Holidays:** No
**Closed Today:** No

#### Branch 2 - Lahug Branch (ID: 4)
**Address:** 3/F Ayala Terraces, Salinas Drive, Lahug, Cebu City, 6000 Cebu
**Phone:** +63 32 234 5679
**Operating Hours:** Mon-Fri: 10:00-18:00, Sat: 10:00-15:00
**Operating Days:** Monday, Tuesday, Wednesday, Thursday, Friday, Saturday
**Open on Holidays:** No
**Closed Today:** No

**Users:**
- **Admins:**
  - mlopez / admin (Miguel Lopez) - miguel.lopez@brightsmile.ph
  - gjose / admin (Grace Jose) - grace.jose@brightsmile.ph
- **Doctors:**
  - fflores / admin (Ferdinand Flores) - ferdinand.flores@brightsmile.ph
  - maquino / admin (Michelle Aquino) - michelle.aquino@brightsmile.ph
- **Secretary:**
  - smendoza / admin (Sofia Mendoza) - sofia.mendoza@brightsmile.ph
- **Assistants:**
  - rramos / admin (Rosa Ramos) - rosa.ramos@brightsmile.ph
  - bcastro / admin (Benjamin Castro) - benjamin.castro@brightsmile.ph

**Patients:** 20 seeded patients with Filipino names, addresses, and comprehensive medical information including:
- Personal details (name, email, phone, date of birth, gender, address)
- Emergency contacts with relationship details
- Medical information (blood type, allergies, medical conditions, current medications)
- Insurance details (PhilHealth, Maxicare, Intellicare)
- Notes for dental care considerations

---

## Summary
- **Total Clinics:** 3 (with taglines)
- **Total Branches:** 4 (with operating hours, days, and holiday settings)
- **Total Users:** 15 (1 super admin + 14 clinic staff)
- **Total Patients:** 40 (20 for SmileCare + 20 for Bright Smile)

All passwords are set to "admin" for testing purposes.
All addresses are based in Cebu City, Philippines with local phone numbers.

## Database Schema Updates

### Clinic and Branch Field Additions
- **Added to Clinics:** `tagline` field for clinic slogans/taglines
- **Added to Branches:**
  - `operating_hours` - JSON string for flexible operating hours (e.g., `{"monday": "9:00-17:00"}`)
  - `operating_days` - JSON array for operating days (e.g., `["monday", "tuesday"]`)
  - `is_open_on_holidays` - Boolean flag for holiday operations
  - `is_closed_today` - Override flag to close branch even when normally open
- **Benefit:** Enhanced clinic and branch management with detailed operating information

### Consent Template Index Fix
- **Fixed:** `idx_consent_templates_code` to be unique per clinic
- **Changed:** From global unique `code` to composite unique `(code, clinic_id)`
- **Benefit:** Same consent template codes can now exist across different clinics
- **Index Name:** `idx_consent_templates_code_clinic`

This ensures that each clinic can have its own set of consent templates with the same codes (e.g., "EXTRACT-001") without conflicts.
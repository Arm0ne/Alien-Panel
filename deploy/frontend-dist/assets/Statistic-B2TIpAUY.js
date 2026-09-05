import{$r as e,In as t,Mi as n,Nn as r,Qr as i,Vi as a,Xn as o,Zn as s,nr as c,zi as l}from"./router-BXCsM_PA.js";import{n as u}from"./index-BCfmTlJD.js";var d=i(`statistic`,[e(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),i(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[e(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[i(`icon`,{verticalAlign:`-0.125em`})]),e(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),e(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[i(`icon`,{verticalAlign:`-0.125em`})])])]),f=l({name:`Statistic`,props:Object.assign(Object.assign({},r.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(e){let{mergedClsPrefixRef:i,inlineThemeDisabled:a,mergedRtlRef:c}=s(e),l=r(`Statistic`,`-statistic`,d,u,e,i),f=t(`Statistic`,c,i),p=n(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:s},common:{cubicBezierEaseInOut:c}}=l.value;return{"--n-bezier":c,"--n-label-font-size":s,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=a?o(`statistic`,void 0,p,e):void 0;return{rtlEnabled:f,mergedClsPrefix:i,cssVars:a?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var e;let{mergedClsPrefix:t,$slots:{default:n,label:r,prefix:i,suffix:o}}=this;return(e=this.onRender)==null||e.call(this),a(`div`,{class:[`${t}-statistic`,this.themeClass,this.rtlEnabled&&`${t}-statistic--rtl`],style:this.cssVars},c(r,e=>a(`div`,{class:`${t}-statistic__label`},this.label||e)),a(`div`,{class:`${t}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},c(i,e=>e&&a(`span`,{class:`${t}-statistic-value__prefix`},e)),this.value===void 0?c(n,e=>e&&a(`span`,{class:`${t}-statistic-value__content`},e)):a(`span`,{class:`${t}-statistic-value__content`},this.value),c(o,e=>e&&a(`span`,{class:`${t}-statistic-value__suffix`},e))))}});export{f as t};